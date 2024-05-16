
'use client';

import Link from 'next/link'
import {Box} from '@chakra-ui/react';

interface Props { }

const NotFound: React.FC<Props> = () => {
    return (
        <Box mt="500h" maxW="800px" m="0 auto" p={5}>
        <Link href="/">Not Found, Return Home</Link>
        </Box>
        
    )
}


export default NotFound;