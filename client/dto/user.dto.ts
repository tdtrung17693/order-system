import { UserRole } from '../constants/user-role'

export interface User {
  id: number
  email: string
  name: string
  role: UserRole
}
