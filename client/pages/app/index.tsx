import React, { useState, useRef, useContext, useEffect } from "react"
import { WebsocketContext } from "../../modules/websocket_provider"
import { useRouter } from "next/router"
import { API_URL, SOCKET_ACTION } from "../../constants"
import autosize from "autosize"
import { AuthContext } from "../../modules/auth_provider"
import ChatBody from "../../components/chat_body"

export type Message = {
  content: string
  client_id: string
  username: string
  room_id: string
  type: "receive" | "self"
  action: (typeof SOCKET_ACTION)[keyof typeof SOCKET_ACTION] // This ensures action must be one of the values (0, 1, 2, or 3)
}

const index = () => {
  const [messages, setMessage] = useState<Array<Message>>([])
  const textarea = useRef<HTMLTextAreaElement>(null)
  const { conn } = useContext(WebsocketContext)
  const [users, setUsers] = useState<Array<{ username: string }>>([])
  const { user } = useContext(AuthContext)

  const router = useRouter()

  useEffect(() => {
    if (conn === null) {
      router.push("/")
      return
    }

    const roomId = conn.url.split("/")[5]
    async function getUsers() {
      try {
        const res = await fetch(`${API_URL}/ws/getClients/${roomId}`, {
          method: "GET",
          headers: { "Content-Type": "application/json" },
        })
        const data = await res.json()

        setUsers(data)
      } catch (e) {
        console.error(e)
      }
    }
    getUsers()
  }, [])

  useEffect(() => {
    if (textarea.current) {
      autosize(textarea.current)
    }

    if (conn === null) {
      router.push("/")
      return
    }

    conn.onmessage = (message) => {
      const m: Message = JSON.parse(message.data)
      // console.log('m',m)
      if (m.action === SOCKET_ACTION.JOIN) {
        setUsers([...users, { username: m.username }])
      }

      if (m.action === SOCKET_ACTION.LEFT) {
        const deleteUser = users.filter((user) => user.username != m.username)
        setUsers([...deleteUser])
        setMessage([...messages, m])
        return
      }

      user?.username == m.username ? (m.type = "self") : (m.type = "receive")
      setMessage([...messages, m])
    }

    conn.onclose = () => {}
    conn.onerror = () => {}
    conn.onopen = () => {}
  }, [textarea, messages, conn, users])

  const sendMessage = () => {
    if (!textarea.current?.value) return
    if (conn === null) {
      router.push("/")
      return
    }

    conn.send(textarea.current.value)
    textarea.current.value = ""
  }

  return (
    <>
      <div className="flex flex-col w-full">
        <div className="p-4 md:mx-6 mb-14">
          <ChatBody data={messages} />
        </div>
        <div className="fixed bottom-0 mt-4 w-full">
          <div className="flex md:flex-row px-4 py-2 bg-grey md:mx-4 rounded-md">
            <div className="flex w-full mr-4 rounded-md border border-blue">
              <textarea
                ref={textarea}
                placeholder="type your message here"
                className="w-full h-10 p-2 rounded-md focus:outline-none"
                style={{ resize: "none" }}
              />
            </div>
            <div className="flex items-center">
              <button
                className="p-2 rounded-md bg-blue text-white"
                onClick={sendMessage}
              >
                Send
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

export default index
