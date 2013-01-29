/**
 * Networking portion of the Jungle framework
 * @author Blake Loring
 */

function JungleNet(details, callbacks) {

	//Store the connection details
	this.details = details;

	//Store the callbacks
	this.callbacks = callbacks;

	//Set logged in to false
	this.loggedIn = false;

	//Open a new socket
	this.socket = new WebSocket(details.Address);

	//Store a reference to this so the callbacks can work
	var thisReference = this;

	//Set the sockets connection listener to the onConnect function
	this.socket.onopen = function(evt) {
		thisReference.onConnect(evt);
	}
	//Route all messages to the right place
	this.socket.onmessage = function(evt) {
		thisReference.onMessage(evt);
	}
	//Set the close listener
	this.socket.onclose = function(evt) {
		thisReference.onClose(evt.reason);
	}
}

/**
 * Function to send a message
 */

JungleNet.prototype.Send = function(object) {
	var encoded = $.toJSON(object);
	this.socket.send(encoded);
}
/**
 * Function to process every message received
 */

JungleNet.prototype.onMessage = function(evt) {

	var incoming = $.parseJSON(evt.data);

	console.log(incoming);

	if (incoming.Command == "Disconnect") {

		//Release the socket
		this.Release();

		//Call the onClose function with the reason given by the message
		this.onClose(incoming.Data);

	} else {

		this.callbacks.data(incoming);

	}
}

JungleNet.prototype.onConnect = function(evt) {

	var loginRequest = new Object();

	loginRequest.Email = this.details.Email;
	loginRequest.Password = this.details.Password;

	this.Send(loginRequest);

	var chatMessage = new Object();

	chatMessage.Command = "Chat";
	chatMessage.Data = "BOOM CHICA";

	this.Send(chatMessage);

	this.callbacks.connected();
}
/**
 * onClose binding. Called when the connection has been closed
 */
JungleNet.prototype.onClose = function(reason) {

	if (reason == null || reason.length == 0) {
		reason = "Unknown";
	}

	this.callbacks.failure(reason);
}

JungleNet.prototype.Release = function() {

	if (this.socket != null) {
		this.socket.onopen = null;
		this.socket.onclose = null;
		this.socket.onmessage = null;
		this.socket.close();
		this.socket = null;
	}

}
/**
 * Connection details, stores username password and address
 */
function JungleNetConnectionDetails(email, password, address) {
	this.Email = email;
	this.Password = password;
	this.Address = address;
}
