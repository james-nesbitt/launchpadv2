package k0s

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"path"

	"regexp"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/jellydator/validation"
	// "github.com/jellydator/validation/is"
	"github.com/k0sproject/version"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/host/network"
	"github.com/Mirantis/launchpad/pkg/util/flags"
	"github.com/Mirantis/launchpad/pkg/util/quote"
	"github.com/Mirantis/launchpad/pkg/util/uploadfile"
)

/**
 * K0S has a HostPlugin which allows K0S specific data to be included in the
 * Host block, and also provided a Host specific docker implementation.
 */

const (
	HostRoleK0S = "k0s"
)

var (
	// K0SManagerHostRoles the Host roles accepted for managers.
	K0SManagerHostRoles = []string{"controller"}
)

func init() {
	host.RegisterHostPluginFactory(HostRoleK0S, &HostPluginFactory{})
}

type HostPluginFactory struct {
	ps []*hostPlugin
}

// HostPlugin build a new host plugin
func (mpf *HostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &hostPlugin{
		h:           h,
		Environment: map[string]string{},
	}

	defaults.Set(p)

	mpf.ps = append(mpf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin
func (mpf *HostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &hostPlugin{
		h:           h,
		Environment: map[string]string{},
	}

	if err := d(p); err != nil {
		return p, err
	}

	if err := defaults.Set(p); err != nil {
		return p, err
	}

	mpf.ps = append(mpf.ps, p)

	return p, nil
}

// Get the K0S plugin from a Host
func HostGetK0s(h *host.Host) *hostPlugin {
	hgk0s := h.MatchPlugin(HostRoleK0S)
	if hgk0s == nil {
		return nil
	}

	hk0s, ok := hgk0s.(*hostPlugin)
	if !ok {
		return nil
	}

	return hk0s
}

// hostPlugin
//
// Implements no generic plugin interfaces.
type hostPlugin struct {
	h *host.Host

	Role             string                   `yaml:"role"`
	Reset            bool                     `yaml:"reset,omitempty"`
	DataDir          string                   `yaml:"dataDir,omitempty"`
	Environment      map[string]string        `yaml:"environment,flow,omitempty"`
	UploadBinary     bool                     `yaml:"uploadBinary,omitempty"`
	K0sBinaryPath    string                   `yaml:"k0sBinaryPath,omitempty"`
	K0sDownloadURL   string                   `yaml:"k0sDownloadURL,omitempty"`
	InstallFlags     flags.Flags              `yaml:"installFlags,omitempty"`
	Files            []*uploadfile.UploadFile `yaml:"files,omitempty"`
	OSIDOverride     string                   `yaml:"os,omitempty"`
	HostnameOverride string                   `yaml:"hostname,omitempty"`
	NoTaints         bool                     `yaml:"noTaints,omitempty"`

	UploadBinaryPath string     `yaml:"-"`
	Metadata         metadata   `yaml:"-"`
	configurer       configurer `yaml:"-"`
}

// HostMetadata resolved metadata for host
type metadata struct {
	K0sBinaryVersion  *version.Version
	K0sBinaryTempFile string
	K0sRunningVersion *version.Version
	K0sInstalled      bool
	K0sExistingConfig string
	K0sNewConfig      string
	K0sJoinToken      string
	K0sJoinTokenID    string
	Arch              string
	IsK0sLeader       bool
	Hostname          string
	Ready             bool
	NeedsUpgrade      bool
	MachineID         string
	DryRunFakeLeader  bool
}

// Id uniquely identify the plugin
func (p hostPlugin) Id() string {
	return fmt.Sprintf("%s:%s", p.h.Id(), "k0s")
}

// RoleMatch what host roles does this host plugin act
func (p hostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleK0S:
		return true
	}

	return false
}

// Validate the plugin configuration
func (p *hostPlugin) Validate() error {
	// For rig validation
	v := validator.New()
	if err := v.Struct(p); err != nil {
		return err
	}

	return validation.ValidateStruct(p,
		validation.Field(p.Role, validation.In("controller", "worker", "controller+worker", "single").Error("unknown role "+p.Role)),
		//validation.Field(p.PrivateAddress, is.IP),
		validation.Field(p.Files),
		validation.Field(p.NoTaints, validation.When(p.Role != "controller+worker", validation.NotIn(true).Error("noTaints can only be true for controller+worker role"))),
		validation.Field(p.InstallFlags, validation.Each(validation.By(quote.ValidateBalancedQuotes))),
	)
}

func (p hostPlugin) getConfigurer() configurer {
	if p.configurer == nil {
		p.configurer = &linuxConfigurer{}
	}
	return p.configurer
}

// **** migrated from k0sctl/cluster/host

var k0sForceFlagSince = version.MustConstraint(">= v1.27.4+k0s.0")

func (p *hostPlugin) SetDefaults() {
	// OS/EXEC SPECIFIC :: if p.OSIDOverride != "" {
	//		p.OSVersion = &rigos.Release{ID: p.OSIDOverride}
	//	}

	// OS/EXEC SPECIFIC ::_ = defaults.Set(h.Connection)

	if p.InstallFlags.Get("--single") != "" && p.InstallFlags.GetValue("--single") != "false" && p.Role != "single" {
		slog.Debug(fmt.Sprintf("%s: changed role from '%s' to 'single' because of --single installFlag", p.Id(), p.Role))
		p.Role = "single"
	}
	if p.InstallFlags.Get("--enable-worker") != "" && p.InstallFlags.GetValue("--enable-worker") != "false" && p.Role != "controller+worker" {
		slog.Debug(fmt.Sprintf("%s: changed role from '%s' to 'controller+worker' because of --enable-worker installFlag", p.Id(), p.Role))
		p.Role = "controller+worker"
	}

	if p.InstallFlags.Get("--no-taints") != "" && p.InstallFlags.GetValue("--no-taints") != "false" {
		p.NoTaints = true
	}

	if dd := p.InstallFlags.GetValue("--data-dir"); dd != "" {
		if p.DataDir != "" {
			slog.Debug(fmt.Sprintf("%s: changed dataDir from '%s' to '%s' because of --data-dir installFlag", p.Id(), p.DataDir, dd))
		}
		p.InstallFlags.Delete("--data-dir")
		p.DataDir = dd
	}
}

// K0sJoinTokenPath returns the token file path from install flags or configurer
func (p *hostPlugin) K0sJoinTokenPath() string {
	if path := p.InstallFlags.GetValue("--token-file"); path != "" {
		return path
	}

	return p.getConfigurer().K0sJoinTokenPath()
}

// K0sConfigPath returns the config file path from install flags or configurer
func (p *hostPlugin) K0sConfigPath() string {
	if path := p.InstallFlags.GetValue("--config"); path != "" {
		return path
	}

	if path := p.InstallFlags.GetValue("-c"); path != "" {
		return path
	}

	return p.getConfigurer().K0sConfigPath()
}

// unquote + unescape a string
func unQE(s string) string {
	unq, err := strconv.Unquote(s)
	if err != nil {
		return s
	}

	c := string(s[0])                                           // string was quoted, c now has the quote char
	re := regexp.MustCompile(fmt.Sprintf(`(?:^|[^\\])\\%s`, c)) // replace \" with " (remove escaped quotes inside quoted string)
	return string(re.ReplaceAllString(unq, c))
}

// K0sInstallCommand returns a full command that will install k0s service with necessary flags
func (p *hostPlugin) K0sInstallCommand(ctx context.Context) (string, error) {
	role := p.Role
	fs := p.InstallFlags

	fs.AddOrReplace(fmt.Sprintf("--data-dir=%s", p.K0sDataDir()))

	switch role {
	case "controller+worker":
		role = "controller"
		fs.AddUnlessExist("--enable-worker")
		if p.NoTaints {
			fs.AddUnlessExist("--no-taints")
		}
	case "single":
		role = "controller"
		fs.AddUnlessExist("--single")
	}

	if !p.Metadata.IsK0sLeader {
		fs.AddUnlessExist(fmt.Sprintf(`--token-file "%s"`, p.K0sJoinTokenPath()))
	}

	if p.IsController() {
		fs.AddUnlessExist(fmt.Sprintf(`--config "%s"`, p.K0sConfigPath()))
	}

	if strings.HasSuffix(p.Role, "worker") {
		hn, nerr := network.HostGetNetwork(p.h).Network(ctx)
		if nerr != nil {
			return "", nerr
		}

		var extra flags.Flags
		if old := fs.GetValue("--kubelet-extra-args"); old != "" {
			extra = flags.Flags{unQE(old)}
		}
		// set worker's private address to --node-ip in --extra-kubelet-args if cloud ins't enabled
		enableCloudProvider, err := p.InstallFlags.GetBoolean("--enable-cloud-provider")
		if err != nil {
			return "", fmt.Errorf("--enable-cloud-provider flag is set to invalid value: %s. (%v)", p.InstallFlags.GetValue("--enable-cloud-provider"), err)
		}
		if !enableCloudProvider && hn.PrivateAddress != "" {
			extra.AddUnlessExist(fmt.Sprintf("--node-ip=%s", hn.PrivateAddress))
		}

		if p.HostnameOverride != "" {
			extra.AddOrReplace(fmt.Sprintf("--hostname-override=%s", p.HostnameOverride))
		}
		if extra != nil {
			fs.AddOrReplace(fmt.Sprintf("--kubelet-extra-args=%s", strconv.Quote(extra.Join())))
		}
	}

	if fs.Include("--force") && p.Metadata.K0sBinaryVersion != nil && !k0sForceFlagSince.Check(p.Metadata.K0sBinaryVersion) {
		slog.Warn(fmt.Sprintf("%s: k0s version %s does not support the --force flag, ignoring it", p.Id(), p.Metadata.K0sBinaryVersion))
		fs.Delete("--force")
	}

	return p.getConfigurer().K0sCmdf("install %s %s", role, fs.Join()), nil
}

// K0sBackupCommand returns a full command to be used as run k0s backup
func (p *hostPlugin) K0sBackupCommand(targetDir string) string {
	return p.getConfigurer().K0sCmdf("backup --save-path %s --data-dir %s", shellescape.Quote(targetDir), p.K0sDataDir())
}

// K0sRestoreCommand returns a full command to restore cluster state from a backup
func (p *hostPlugin) K0sRestoreCommand(backupfile string) string {
	return p.getConfigurer().K0sCmdf("restore --data-dir=%s %s", p.K0sDataDir(), shellescape.Quote(backupfile))
}

// IsController returns true for controller and controller+worker roles
func (p *hostPlugin) IsController() bool {
	return p.Role == "controller" || p.Role == "controller+worker" || p.Role == "single"
}

// K0sServiceName returns correct service name
func (p *hostPlugin) K0sServiceName() string {
	switch p.Role {
	case "controller", "controller+worker", "single":
		return "k0scontroller"
	default:
		return "k0sworker"
	}
}

func (p *hostPlugin) k0sBinaryPathDir() string {
	return path.Dir(p.getConfigurer().K0sBinaryPath())
}

// K0sDataDir returns the data dir for the host either from host.DataDir or the default from configurer's DataDirDefaultPath
func (p *hostPlugin) K0sDataDir() string {
	if p.DataDir == "" {
		return p.getConfigurer().DataDirDefaultPath()
	}
	return p.DataDir
}

// CheckHTTPStatus will perform a web request to the url and return an error if the http status is not the expected
func (p *hostPlugin) CheckHTTPStatus(ctx context.Context, url string, expected ...int) error {
	status, err := p.getConfigurer().HTTPStatus(ctx, p.h, url)
	if err != nil {
		return err
	}

	for _, e := range expected {
		if status == e {
			return nil
		}
	}

	return fmt.Errorf("expected response code %d but received %d", expected, status)
}

// ExpandTokens expands percent-sign prefixed tokens in a string, mainly for the download URLs.
// The supported tokens are:
//
//   - %% - literal %
//   - %p - host architecture (arm, arm64, amd64)
//   - %v - k0s version (v1.21.0+k0s.0)
//   - %x - k0s binary extension (.exe on Windows)
//
// Any unknown token is output as-is with the leading % included.
func (p *hostPlugin) ExpandTokens(ctx context.Context, input string, k0sVersion *version.Version) string {
	if input == "" {
		return ""
	}
	builder := strings.Builder{}
	var inPercent bool
	for i := 0; i < len(input); i++ {
		currCh := input[i]
		if inPercent {
			inPercent = false
			switch currCh {
			case '%':
				// Literal %.
				builder.WriteByte('%')
			case 'p':
				// Host architecture (arm, arm64, amd64).
				builder.WriteString(p.Metadata.Arch)
			case 'v':
				// K0s version (v1.21.0+k0s.0)
				builder.WriteString(url.QueryEscape(k0sVersion.String()))
			case 'x':
				// K0s binary extension (.exe on Windows).
				if exec.HostGetPlatform(p.h).IsWindows(ctx) {
					builder.WriteString(".exe")
				}
			default:
				// Unknown token, just output it with the leading %.
				builder.WriteByte('%')
				builder.WriteByte(currCh)
			}
		} else if currCh == '%' {
			inPercent = true
		} else {
			builder.WriteByte(currCh)
		}
	}
	if inPercent {
		// Trailing %.
		builder.WriteByte('%')
	}
	return builder.String()
}
