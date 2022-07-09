export const ItemsPerPage = process.env.NEXT_PUBLIC_ITEMS_PER_PAGE
    ? parseInt(process.env.NEXT_PUBLIC_ITEMS_PER_PAGE, 10) 
    : 1
