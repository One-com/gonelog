# Deprecated

This package was an experiment if it was possible to make a context object also a Logger. The conclusion is that it'll never work due to the nature of the (now stdlib) context package, where derived contexts are created with package global functions. Any extention of context will be hidden from use when a new context is derived.

Please do not use this package.

