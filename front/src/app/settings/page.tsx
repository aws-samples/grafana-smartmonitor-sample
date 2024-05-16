'use client';

import React, { useContext, useEffect, useMemo, useState } from 'react';
import { AlertDialog, AlertDialogBody, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay, Box,Button,Checkbox,Flex,Input,Table,Tbody,Td,Text, Th, Thead, Tr, useDisclosure, useToast } from '@chakra-ui/react';
import { AuthContext } from '@components/AuthContext';
import { useRouter } from 'next/navigation';
import { useTable } from 'react-table';
import ConnationActionCell from '@components/ConnationActionCell';

import { API_URL } from '../../constants';
import { ColorModeSwitch } from '@components/ColorModeSwitch';


const SettingsPage = () => {
  const [data, setData] = useState([]);

  const [showDetail,setShowDetail] = useState(false)

  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = React.useRef();
  const { isAuthenticated } = useContext(AuthContext);
  const [config, setConfig] = useState(null);
  const router = useRouter()


  const [confirmIsOpen, setConfirmOpen] = useState(false);

  const [newConnection, setNewConnection] = useState({
    conn_name: '',
    conn_username: '',
    conn_password:'',
    conn_url: '',
  });


  const [editConn, setEditConn] = useState(null);
  const [deleteConn, setDeleteConn] = useState(null);


  const toast = useToast();

  const showWarning = () => {
    toast({
      position: "top",
      render: () => (
        <Box
          color="white"
          p={3}
          bg="orange.500"
          borderRadius="md"
          boxShadow="md"
        >
          <Text fontWeight="bold">Warning</Text>
          <Text>Please check conneciton name !</Text>
        </Box>
      ),
    });
  };
  

  const fetchConnection = async () => {
    
    try {
      const token = localStorage.getItem('token');
      // Construct the headers object with the Authorization header
      const headers = {
        'Authorization': `Bearer ${token}`,
        // Add any other headers if needed
      };
      const response = await fetch(`${API_URL}/connections`,{headers});
      const data = await response.json();
     
      if (data.connections) {
        setData(data.connections);
      }

    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      
    }
  };
  
  const fetchConfig = async () => {
    try {
      const token = localStorage.getItem('token');
      // Construct the headers object with the Authorization header
      const headers = {
        'Authorization': `Bearer ${token}`,
        // Add any other headers if needed
      };

      const response = await fetch(`${API_URL}/config`,{
        headers,
      });
      const data = await response.json();
      setConfig(data);
    } catch (error) {
      console.error('Error fetching config:', error);
    }
  };

  



  useEffect(() => {
    fetchConfig();
    fetchConnection()

    
  }, []);

  const handleInputChange = (e) => {
    setNewConnection({ ...newConnection, [e.target.name]: e.target.value });
  };


  const handleAddConn = async () => {
    try {
      const token = localStorage.getItem('token');
      
      const response = await fetch(`${API_URL}/connection`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(newConnection),
      });

      if (response.ok) {
        // Reset the form and fetch the updated data
        setNewConnection({
          conn_name: '',
          conn_url: '',
          conn_username: '',
          conn_password: '',
        });
        fetchConnection();
      } else {
        showWarning();
      }
    } catch (error) {
      console.error('Error adding job:', error);
    } finally {
      onClose();
    }
  };

  const handleUpdateConn = async () => {
    console.log(editConn)

    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`${API_URL}/connection/${editConn.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(editConn),
      });

      if (response.ok) {
        // Fetch the updated data
        fetchConnection();
      } else {
        console.error('Error updating job:', response.status);
      }
    } catch (error) {
      console.error('Error updating job:', error);
    } finally {
      setEditConn(null);
      onClose();
    }
  }

  const handleDeleteConn = async (id) => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`${API_URL}/connection/${id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        // Fetch the updated data
        fetchConnection();
      } else {
        console.error('Error deleting job:', response.status);
      }
    } catch (error) {
      console.error('Error deleting job:', error);
    } finally {
      setConfirmOpen(false)
    }
  }

 

  const columns = useMemo(
    () => [
      { Header: 'ID', accessor: 'id' },
      { Header: 'Connection Name', accessor: 'conn_name' },
      { Header: 'Username', accessor: 'conn_username' },
      { Header: 'URL', accessor: 'conn_url' },

      {
        Header: 'ACTION',
        Cell: ({ row }) => <ConnationActionCell
          row={row}
          setEditItem={setEditConn}
          onOpen={onOpen}
          setConfirmOpen={setConfirmOpen}
          setDeleteItem={setDeleteConn}
        />,
      },
    ],
    []
  );

  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable({ columns, data });
  
  return  (
    <>
    <Box mt="500h" maxW="800px" m="0 auto" p={5}>
      <Text fontSize="2xl" fontWeight="bold" mb={4}>
        Settings
      </Text>
      
     <ColorModeSwitch />
      
      <Button mb={2} onClick={()=>setShowDetail(!showDetail)}>Detail </Button>


      {config ? (
        showDetail&&(<pre>{JSON.stringify(config, null, 2)}</pre>)
      ) : (
        <Text>Loading configuration...</Text>
      )}
    <hr/>
       <Table {...getTableProps()}>
          <Thead>
            {headerGroups.map((headerGroup) => {
              const { key, ...restHeaderGroupProps } =
                headerGroup.getHeaderGroupProps();
              return (
                <Tr key={key} {...restHeaderGroupProps}>
                  {headerGroup.headers.map((column) => {
                    const { key, ...restColumn } = column.getHeaderProps();
                    return (
                      <Th key={key} {...restColumn}>
                        {column.render("Header")}
                      </Th>
                    );
                  })}
                </Tr>
              );
            })}
          </Thead>
          <Tbody {...getTableBodyProps}>
            {rows.map((row) => {
              prepareRow(row);
              const { key, ...restRowProps } = row.getRowProps();
              return (
                <Tr key={key} {...restRowProps}>
                  {row.cells.map((cell) => {
                    const { key, ...restCellProps } = cell.getCellProps();
                    return (
                      <Td key={key} {...restCellProps}>
                        {cell.render("Cell")}
                      </Td>
                    );
                  })}
                </Tr>
              );
            })}
          </Tbody >
        </Table>
   <Flex justifyContent="center" mt={4} mb={4}>
    <Button mt={2} onClick={() => {
          console.log("Click add connections")
          setEditConn(null)
          onOpen()
        }
        }>Add Grafana Connect</Button></Flex>
  </Box>
  {/*Confrim Alert */}
  <AlertDialog isOpen={confirmIsOpen} leastDestructiveRef={cancelRef} onClose={() => { setConfirmOpen(false) }}>
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Connection
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to delete this connection? <hr />
              {deleteConn?.id}, {deleteConn?.conn_name}, {deleteConn?.conn_url}
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={() => { setConfirmOpen(false) }}>
                Cancel
              </Button>
              <Button colorScheme="red" ml={3} onClick={() => handleDeleteConn(deleteConn?.id)}>
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>

  <AlertDialog isOpen={isOpen} leastDestructiveRef={cancelRef} onClose={onClose}>
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              {editConn ? 'Edit Connection' : 'Add Connection'}
            </AlertDialogHeader>
            <AlertDialogBody>
              Connection Name<Input
                name="conn_name"
                placeholder="Connection Name"
                value={editConn ? editConn.conn_name : newConnection.conn_name}
                onChange={(e) => {
                  if (editConn) {
                    setEditConn({ ...editConn, conn_name: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              />

              Grafana URL
              <Input
                name="conn_url"
                placeholder="Connection URL"
                value={editConn ? editConn.conn_url : newConnection.conn_url}
                onChange={(e) => {
                  if (editConn) {
                    setEditConn({ ...editConn, conn_url: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              />
              Grafana Admin User
              <Input
                name="conn_username"
                placeholder="Connection Admin User"
                value={editConn ? editConn.conn_username : newConnection.conn_username}
                onChange={(e) => {
                  
                  if (editConn) {
                    setEditConn({ ...editConn, conn_username: e.target.value });
                  } else {
                     handleInputChange(e);
                  }
                }}
                mb={2}
              />
              Grafana Admin Password
              <Input
                type="password"
                name="conn_password"
                placeholder="Connection Admin Pasword"
                value={editConn ? editConn.conn_password : newConnection.conn_password}
                onChange={(e) => {
                  if (editConn) {
                    setEditConn({ ...editConn, conn_password: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              />
              
              
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              {editConn ? (
                <Button colorScheme="blue" onClick={handleUpdateConn} ml={3}>
                  Update
                </Button>
              ) : (
                <Button colorScheme="blue" onClick={handleAddConn} ml={3}>
                  Add
                </Button>
              )}
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
  </>
  )
  ;
};

export default SettingsPage;

