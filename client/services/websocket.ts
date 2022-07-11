import auth from "./auth"

let wsInstance: WebSocket

function buildWsEndpoint(jwtToken: string) {
    return `${process.env.NEXT_PUBLIC_WEBSOCKET_ENDPOINT}?jwt=${jwtToken}`
}

let onMessageFns: (() => void)[] = []
export const websocket = {
    initialized: false,
    init() {
        if (!auth.authenticated) return
        if (this.initialized) return 
        wsInstance = new WebSocket(buildWsEndpoint(auth.getAccessToken()!))
        
        wsInstance.addEventListener("message", event => {
            console.log("message from server: ", )
            onMessageFns.forEach(fn => fn())
        })
        this.initialized = true
    },
    destroy() {
        onMessageFns = []
        wsInstance.close()
    },
    onMessage(fn: () => void) {
        onMessageFns.push(fn)
    }
}
