"use client";

import { useContext, useEffect } from 'react';
import { AuthContext } from '@components/AuthContext';
import { useRouter } from 'next/navigation';

const LogoutPage = () => {
const { logout } = useContext(AuthContext);
  const router = useRouter();

  useEffect(() => {
    // Clear authentication data from localStorage
    logout()
    // Redirect to the login page
    router.push('/login');
  }, [router]);

  return <div>Logging out...</div>;
};

export default LogoutPage;
