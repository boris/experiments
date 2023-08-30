// pragma: Specifies the version of Solidity compiler to use (from `truffle version`)
pragma solidity ^0.8.0;

// contract: Defines the contract name to be used in the blockchain
contract HelloWorld{
    // private: Makes the variable private, only available to the contract
    string private hello_msg = "Hello World!";

    // function: Defines a function to be used in the contract
    // public: Makes the function public, available to the outside world
    // view: Makes the function a read-only function, does not modify the state
    // of the contract
    function getHelloMsg() public view returns (string memory) {
        return hello_msg;
    }
}
