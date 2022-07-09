import { UserSignUp, UserTokenResponse } from '../dto/auth.dto'
import { http } from './http'

export const users = {
  register(data: UserSignUp): Promise<UserTokenResponse> {
    return http
      .post<UserTokenResponse>('/register', data)
      .then((response) => response.data)
  },
}
