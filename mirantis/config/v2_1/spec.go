package v21

/**
 * V2.1 Launchpad Spec
 *
 * The 2.1 spec is a pull move to component based approach, dropping
 * the separates hosts declaration, and considering all of the products
 * as components.
 * This makes all components optional, and allows complete freedom for
 * component injection, relying on dependencies to enforce coupling
 * validation.
 *
 * ------------------------------
 * components:
 *   hosts:
 *     one:
 *       mcr:
 *         role: manager
 *     two:
 *       rig:
 *         ssh:
 *           address: 127.0.0.1
 *           key_path: ./key.pem
 *           user: sshuser
 *
 *   k0s:
 *     version: 1.24
 *   registry
 *     handler: msr3
 *     version: 3.1.1
 * ------------------------------
 *
 */

// Spec defines projec spec.
type Spec struct {
	Components SpecComponents `yaml:"components" validate:"required,min=1"`
}
