// Drawing
const canvas = document.getElementById("canvas")
const clearCanvas = document.getElementById("clearCanvas")

// Chat
const chatBox = document.getElementById("chatBox")
const chatMessage = document.getElementById("chatMessage")
const sendChat = document.getElementById("sendChat")

const wsPath = location.origin.replace(/^http/, 'ws') + "/ws"

class Server {
  constructor() {
    this.ws = new WebSocket(wsPath)
    this.types = new Map()

    const self = this
    this.ws.onmessage = function(event) {
      const msg = JSON.parse(event.data)
      if (!self.types.has(msg.Type)) {
        console.log("undifined type:", msg.Type)
        return
      }
      const handler = self.types.get(msg.Type)
      handler(msg)
    }
  }
  addHandler(type, dataHandler) {
    this.types.set(type, dataHandler)
  }
  send(type, data) {
    this.ws.send(JSON.stringify({ Type: type, Data: data }))
  }
  restart() {
    this.ws = new WebSocket(wsPath)
  }
}
const srv = new Server ()
srv.addHandler("chat", (data) => {
  const p = document.createElement("p")
  p.textContent = data.username + ": " + data.message
  chatBox.appendChild(p)
})

class DrawingBoard {
  constructor() {
    this.ctx = canvas.getContext("2d")
    this.drawing = false

    const self = this
    canvas.addEventListener("mousedown", function(event) {
      self.drawing = true
      self.ctx.beginPath()
      self.ctx.moveTo(event.clientX - canvas.offsetLeft, event.clientY - canvas.offsetTop)
    })
    canvas.addEventListener("mousemove", function(event) {
      if (self.drawing) {
        self.ctx.lineTo(event.clientX - canvas.offsetLeft, event.clientY - canvas.offsetTop)
        self.ctx.stroke()
      }
    })
    canvas.addEventListener("mouseup", function(event) {
      if (self.drawing) {
        srv.send("draw", canvas.toDataURL())
        self.drawing = false
      }
    })
  }
  clear() {
    this.ctx.clearRect(0, 0, canvas.width, canvas.height)
  }
}
const board = new DrawingBoard()

clearCanvas.addEventListener("click", function(event) {
  event.preventDefault()
  
  board.clear()
})

sendChat.addEventListener("click", function(event) {
  event.preventDefault()
  chatMessage.value.trim()
  if (0 < chatMessage.value) {
    return
  }
  srv.send("chat", chatMessage.value)
  chatMessage.value = ""
})
