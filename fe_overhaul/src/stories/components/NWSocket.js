
const versionNum = '1.0.0';
const versionTag = 'NodeWars:' + versionNum;
const confirmationPhrase = 'Welcome to NodeWars';


class NWSocket {
	constructor(parser, debug) {
		this.debug = debug || false
		let ws_protocol = 'ws://'
		if (location.protocol == 'https:')
			ws_protocol = 'wss://'

		// const ws = new WebSocket(ws_protocol + window.location.host + '/ws');
		this.ws = new WebSocket(ws_protocol + 'localhost:8080' + '/ws');


		this.ws.addEventListener('open', () => {
					// Try to handshake
					if (debug) console.log('<NWSocket> >', versionTag)
					this.ws.send(versionTag)
				});

		this.ws.addEventListener('message', (e)=>{
			if (e.data == confirmationPhrase) {
				if (debug) console.log('<NWSocket> > Handshake succesful <')

				// turn on normal message parsing
				this.ws.addEventListener('message', parser.handle);

				// handle server terminating connection
				this.ws.addEventListener('close', (d) => {console.log("server severed connection:", d)});
				return
			}
			console.log('Server said:', e.data)
			throw "Error: failed to negotiate handshake with server"
			this.ws.close()
		}, { once: true });
	}

	send(msg) {
		if (this.debug) console.log('<NWSocket> Sending message,',msg)
		this.ws.send(msg)
	}
}

export default NWSocket