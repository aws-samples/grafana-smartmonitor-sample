"use client";

import React, {  useContext, useEffect, useState } from 'react'
import {
  Box,
  CloseButton,
  Flex,
  Icon,
  useColorModeValue,
  Text,
  useDisclosure,
  BoxProps,
  FlexProps,
} from '@chakra-ui/react'
import {
  FiHome,
  FiTrendingUp,
  FiSettings,
  FiLogOut ,
  FiLogIn,
} from 'react-icons/fi'
import { IconType } from 'react-icons'
import {MdOutlineSchedule } from "react-icons/md";

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { AuthContext } from './AuthContext';

const isBrowser = typeof window !== 'undefined';


interface LinkItemProps {
  name: string;
  icon: React.ComponentType<any>;
  href: string;
  protect: boolean;
}



interface NavItemProps extends FlexProps {
  icon: IconType;
  href: string;
 
  children: React.ReactNode;
}
const LinkItems: Array<LinkItemProps> = [
  { name: 'Home', icon: FiHome, href: '/',protect:false },
  { name: 'Metrics', icon: FiTrendingUp, href: '/metrics',protect:true},
  { name: 'Scheduler', icon: MdOutlineSchedule, href: '/scheduler' ,protect:true},
  { name: 'Settings', icon: FiSettings, href: '/settings',protect:true },
];

const Sidebar = () => {
  const { onClose } = useDisclosure();

  return (
    <Box minH="100vh" minW="10vw" bg={useColorModeValue('gray.100', 'gray.900')}>
      
      <SidebarContent onClose={() => onClose} display={{ base: 'none', md: 'block' }} />
    </Box>
  );
};


interface SidebarProps extends BoxProps {
  onClose: () => void
}


const SidebarContent = ({ onClose, ...rest }: SidebarProps) => {
  
  const { isAuthenticated } = useContext(AuthContext);


  return (
    <Box
      bg={useColorModeValue('white', 'gray.900')}
      borderRight="1px"
      borderRightColor={useColorModeValue('gray.200', 'gray.700')}
      w={{ base: 'full', md: 60 }}
      pos="fixed"
      h="full"
      {...rest}>
      <Flex h="20" alignItems="center" mx="8" justifyContent="space-between">
      <Text fontSize="2xl" fontFamily="monospace" fontWeight="bold">
          Grafana<br/>SmartMonitor
        </Text>
        <Text fontSize="2xl" fontFamily="monospace" fontWeight="bold">
          
        </Text>
        <CloseButton display={{ base: 'flex', md: 'none' }} onClick={onClose} />
        
       
      </Flex>
      
      {isBrowser && (
        isAuthenticated
          ? LinkItems.map((link, idx) => (
              <NavItem key={link.name} icon={link.icon as IconType} href={link.href}>
                {link.name}
              </NavItem>
            ))
          : LinkItems.filter((link) => !link.protect).map((link, idx) => (
              <NavItem key={link.name} icon={link.icon as IconType} href={link.href}>
                {link.name}
              </NavItem>
            ))
      )}

      {isBrowser && (
        isAuthenticated===false&&<NavItem key={"login"} icon={FiLogIn as IconType} href={'/login'}>
        Login
      </NavItem>
      )}

      {isBrowser && (
        isAuthenticated===true&&<NavItem key={"logout"} icon={FiLogOut as IconType} href={'/logout'}>
        Logout
      </NavItem>
      )}
      
    </Box>
  )
}


const NavItem = ({ icon, children, href, ...rest }: NavItemProps) => {

  const [isClient, setIsClient] = useState(false);
  const router = useRouter();
  const { isAuthenticated } = useContext(AuthContext);

  const protect_href=isAuthenticated?href:href=="/"?href:'/login'
  
  // useEffect(() => {
  //   if (!isAuthenticated && href !== '/login') {
  //     router.push('/login');
  //   }
  // }, [isAuthenticated, href, router]);

  useEffect(() => {
    setIsClient(true)
    
    
}, [])

if (!isClient) return 
  return (
    <>
    <div>
    {isClient&&<Link href={protect_href} passHref>
      <Box
        style={{ textDecoration: 'none' }}
        _focus={{ boxShadow: 'none' }}
      >
        <Flex
          align="center"
          p="4"
          mx="4"
          borderRadius="lg"
          role="group"
          cursor="pointer"
          _hover={{
            bg: 'cyan.400',
            color: 'white',
          }}
          {...rest}
        >
          {icon && (
            <Icon
              mr="4"
              fontSize="16"
              _groupHover={{
                color: 'white',
              }}
              as={icon}
            />
          )}
          {children}
        </Flex>
      </Box>
    </Link>}
    </div>
    </>
  );
};


export default Sidebar;
