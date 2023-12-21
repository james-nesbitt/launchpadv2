package product

type Product interface {
	// Debug output state to log and return a data structure that can be used to debug the product for tracing
	Debug() interface{}
}
