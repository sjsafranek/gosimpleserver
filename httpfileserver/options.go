package httpfileserver

// // Option is the type all options need to adhere to
// type Option func(server *FileServer)

// // OptionNoCache disables the caching
// func OptionNoCache(disable bool) Option {
// 	return func(server *FileServer) {
// 		server.optionDisableCache = disable
// 	}
// }

// // OptionMaxBytes sets the maximum number of bytes per file to cache,
// // the default is 10 MB
// func OptionMaxBytes(optionMaxBytesPerFile int) Option {
// 	return func(server *FileServer) {
// 		server.optionMaxBytesPerFile = optionMaxBytesPerFile
// 	}
// }
