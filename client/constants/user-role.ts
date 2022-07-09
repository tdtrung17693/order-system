export enum UserRole {
  User,
  Vendor,
}

export function parseRole(roleString: string) {
    const val = parseInt(roleString, 10)
    if (isNaN(val) || [UserRole.User, UserRole.Vendor].indexOf(val) < 0) return UserRole.User
    
    return UserRole.Vendor
}
