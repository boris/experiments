# Considerations:
- Truffle must be running and the port must be the same as the one defined on
    `truffle-config.js:L68L71`
- Compiler version might need to be updated on `truffle-config.js:L110`

# Commands to be used:

```
truffle compile
truffle deploy
truffle console
```

Once in the truffle console;

```
truffle(development)> HelloWorld.deployed().then(function(instance) {return instance});
TruffleContract {
(...)
}
truffle(development)> HelloWorld.deployed().then(function(instance) {return instance.getHelloMsg()});
'Hello World!'
```
