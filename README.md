[tetratelabs/proxy-wasm-go-sdk](https://github.com/tetratelabs/proxy-wasm-go-sdk) has been archived due to memory leaks caused by the built-in garbage collector (gc) in TinyGo. 

This repository is a stable alternative, where Higress has replaced TinyGo's built-in gc with [bdwgc](https://github.com/ivmai/bdwgc) (being used by large projects like Unity3D), largely avoiding the memory leak issue (though it still exists in extreme cases, such as: https://github.com/wasilibs/nottinygc/issues/46). 

The Higress community has [40+ wasm plugins](https://higress.io/en/plugin/), and none of them have memory leak issues.

Furthermore, Higress continues to optimize through [nottinygc](https://github.com/higress-group/nottinygc), which theoretically can prevent the memory leaks in those extreme scenarios.
