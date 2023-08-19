const CONTRACT_ADDRESS = "0x5aa97C24D90f309F5CdeF8741d34f1A1F231d878"
const META_DATA_URL = "ipfs://bafyreicxmahuwolxapxuwcdsgfzzhbmp7pfomxwrcuyk3h6llojcdediam"

async function mintNFT(contractAddress, metaDataURL) {
   const ExampleNFT = await ethers.getContractFactory("ExampleNFT")
   const [owner] = await ethers.getSigners()
   await ExampleNFT.attach(contractAddress).mintNFT(owner.address, metaDataURL)
   console.log("NFT minted to: ", owner.address)
}

mintNFT(CONTRACT_ADDRESS, META_DATA_URL)
   .then(() => process.exit(0))
   .catch((error) => {
       console.error(error);
       process.exit(1);
   });
