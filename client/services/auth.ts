import Router from 'next/router'
import { UserRole } from '../constants/user-role'
import { User } from '../dto/user.dto'
import { Maybe } from '../types/maybe'
import { http } from './http'

const ACCESS_TOKEN_KEY = 'access_token'
const stubUser: Maybe<User> = null

interface AuthService {
  authenticated: boolean
  initialized: boolean
  user: Maybe<User>
  getAccessToken: () => Maybe<string>
  setAccessToken: (token: string) => void
  onLogout: (fn: () => void) => void
  onLogin: (fn: () => void) => void
  init: () => Promise<Maybe<User>>
  login: (email: string, password: string) => Promise<any>
  logout: () => void
}
const onLogoutFns: (() => void)[] = []
const onLogInFns: (() => void)[] = []
const auth: AuthService = {
  authenticated: false,
  initialized: false,
  user: stubUser,
  onLogout(fn) {
    onLogoutFns.push(fn)
  },
  onLogin(fn) {
    onLogInFns.push(fn)
  },
  async login(email: string, password: string) {
    let response = await http.post('/login', {
      email,
      password,
    })
    localStorage.setItem(ACCESS_TOKEN_KEY, response.data.accessToken)
    response = await http.get('/me')
    this.user = response.data
    this.authenticated = true
    onLogInFns.forEach((fn) => fn())
  },
  async init() {
    if (this.initialized) return
    try {
      const token = this.getAccessToken()

      if (!token) return

      const response = await http.get('/me')
      this.user = response.data
      this.authenticated = true

      return this.user
    } catch (error) {
      this.logout()
    } finally {
      this.initialized = true
    }
  },
  setAccessToken(token: string) {
    localStorage.setItem(ACCESS_TOKEN_KEY, token)
  },
  getAccessToken() {
    return localStorage.getItem(ACCESS_TOKEN_KEY)
  },
  logout() {
    this.authenticated = false
    this.user = null
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    onLogoutFns.forEach((fn) => fn())
    Router.push('/')
  },
}

export default auth
