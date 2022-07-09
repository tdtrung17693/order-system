import { UserRole } from '../constants/user-role'

export interface UserSignUp {
  name: string
  email: string
  password: string
  confirmPassword: string
  role: UserRole
}

export interface UserSignIn {
  email: string
  password: string
}

export interface UserTokenResponse {
  accessToken: string
}
