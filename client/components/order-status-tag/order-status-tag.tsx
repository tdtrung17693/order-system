import { Tag } from 'antd'
import { OrderStatus } from 'constants/order'

interface OrderStatusTagProps {
  status: OrderStatus
}
export const OrderStatusTag: React.FC<OrderStatusTagProps> = (props) => {
  const status = props.status
  if (status === OrderStatus.Cancelled) {
    return <Tag color="volcano">{status}</Tag>
  } else if (status === OrderStatus.Shipped) {
    return <Tag color="green">{status}</Tag>
  }
  return <Tag color="geekblue">{status}</Tag>
}
