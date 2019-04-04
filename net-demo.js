const net = require('net')
const { StringDecoder } = require('string_decoder')
const decoder = new StringDecoder('utf8')

let client = net.connect('8899', 'localhost', () => {
	console.log('connect success')
})

client.on('data', data => {
	console.log('receive', decoder.write(data))
})

// var message = { path: 'user/login', data: '{"username":"admin","userpwd":"123456"}' }
// client.write(Buffer.from(JSON.stringify(message)))

// var _sms = JSON.stringify({ source: { username: 'admin' }, data: 'hello' })
// var sms = { path: 'sms/add', data: _sms }
// client.write(Buffer.from(JSON.stringify(sms)))

var users = [
	'{"username":"admin", "userpwd":"123456","usernickname":"admin"}',
	'{"username":"ceshi1", "userpwd":"123456","usernickname":"ceshi1"}',
	'{"username":"ceshi2", "userpwd":"123456","usernickname":"ceshi2"}'
]

async function addUser() {
	for (let i = 0; i < users.length; i++) {
		console.log(Date.now())
		await (() => {
			return new Promise((resolve, reject) => {
				setTimeout(() => {
					resolve()
				}, 1000)
			})
		})()
		user = users[i]
		console.log('send:', user)
		client.write(Buffer.from(JSON.stringify({ path: 'user/add', data: user })))
	}
}
addUser()
