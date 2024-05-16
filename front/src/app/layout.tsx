'use client';

import { RecoilRoot } from 'recoil'
import { Box, ChakraProvider, Flex } from '@chakra-ui/react'
import Sidebar from '../components/Sidebar';
import { AuthProvider } from '@components/AuthContext';


export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {

  return (
    <html lang="en">
      <body>
        <RecoilRoot>
          <AuthProvider>
            <ChakraProvider >
            
              <Flex h="100vh">
                <Sidebar />
                <Box flex="1" p={6}>
                  {children}
                </Box>
              </Flex>
             
            </ChakraProvider>
          </AuthProvider>
        </RecoilRoot>
      </body>
    </html>
  )
}
