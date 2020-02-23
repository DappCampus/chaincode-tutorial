const hfc = require('fabric-client');
const path = require('path');

const NETWORK_CONNECTION_PROFILE_PATH = path.join(__dirname, 'network-config.yaml')
const ORG1_CONNECTION_PROFILE_PATH = path.join(__dirname, 'org1.yaml')

// Org & User
const USER_NAME = 'my'
const MSP_ID = 'Org1MSP'
const PEER_NAME = 'peer0.org1.example.com'
const CHANNEL_NAME = 'mychannel'

const CHAINCODE_ID = 'erc20-transfer'
const CHAINCODE_EVENT = 'transferEvent'

const CRYPTO_CONTENT = {
    privateKey: '/Users/kyung/.fabric-vscode/environments/Local Fabric/wallets/Org1/admin/keystore/2ebbd4c70874a0a7b017894df7b9ca7bd29d5b0dc0846484f455fb2b8a4fbedb_sk',
    signedCert: '/Users/kyung/.fabric-vscode/environments/Local Fabric/wallets/Org1/admin/signcerts/cert.pem'
  }

async function subscribeEvent() {
  try {
    const client = await getClient()
    const channel = await client.getChannel(CHANNEL_NAME)

    let eventHub = channel.newChannelEventHub(PEER_NAME);

    let chaincodeListener = await eventHub.registerChaincodeEvent(CHAINCODE_ID, CHAINCODE_EVENT,

            // onEvent
            (chaincodeEvent)=>{
                console.log(`chaincode event emiited: ${chaincodeEvent.chaincode_id}  ${chaincodeEvent.event_name}  ${new String(chaincodeEvent.payload)}`)
            },
            // onError
            (err)=>{
                console.log('chaincode event error: ', err)
            }
        )

        await eventHub.connect(true)
        console.log('chaincodeEvenrHandler started with handler_id=',chaincodeListener)

  } catch (e) {
    console.log(`error: ${e}`)
    process.exit(1)
  }
}

const getClient = async () => {

    // setup the instance
    const client = hfc.loadFromConfig(NETWORK_CONNECTION_PROFILE_PATH)

    // Call the function for initializing the credentials store on file system
    client.loadFromConfig(ORG1_CONNECTION_PROFILE_PATH)
    await client.initCredentialStores()

    let opts = { username: USER_NAME, mspid: MSP_ID, cryptoContent: CRYPTO_CONTENT, skipPersistence: true }
    let user = await client.createUser(opts)
    await client.setUserContext(user)

    return client
}

subscribeEvent()
