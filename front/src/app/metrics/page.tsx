'use client'
import React, { useContext, useEffect, useMemo, useState } from 'react';
import { Link, Table, Thead, Tbody, Tr, Th, Td, Input, Button, Box, Image, Flex, Spinner, AlertDialog, AlertDialogBody, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay, useDisclosure, Icon, Select } from '@chakra-ui/react';
import { useTable } from 'react-table';


import ActionCell from "@components/ActionCell"


import { API_URL } from '../../constants';


interface Connection {
  conn_name: string;
}

const TableComponent = () => {

  const [data, setData] = useState([]);
  const [connections,setConnections] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  const { isOpen, onOpen, onClose } = useDisclosure();

  const [confirmIsOpen, setConfirmOpen] = useState(false);

  const cancelRef = React.useRef();



  const [newItem, setNewItem] = useState({
    project: '',
    catalog: '',
    item_desc: '',
    item_condition: '',
    conn_name:'',
    dashboard_url: '',
  });


  const [editItem, setEditItem] = useState(null);
  const [deleteItem, setDeleteItem] = useState(null);


  const fetchData = async () => {
    console.log(fetchData)
    setIsLoading(true);
    try {
      const token = localStorage.getItem('token');
      // Construct the headers object with the Authorization header
      const headers = {
        'Authorization': `Bearer ${token}`,
        // Add any other headers if needed
      };

      const response = await fetch(`${API_URL}/metrics`,{
        headers,
      });
      const data = await response.json();
      if(data.metrics){
        setData(data.metrics);
      }
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setTimeout(() => setIsLoading(false), 1000);
    }
  };


  const fetchConn = async () => {
    setIsLoading(true);
    try {
      const token = localStorage.getItem('token');
      // Construct the headers object with the Authorization header
      const headers = {
        'Authorization': `Bearer ${token}`,
        // Add any other headers if needed
      };

      const response = await fetch(`${API_URL}/connections/name`,{
        headers,
      });
      const data = await response.json();
      if(data.connections){
        console.log(data.connections.map((conn: Connection) => conn.conn_name))
        setConnections(data.connections);
      }
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setTimeout(() => setIsLoading(false), 1000);
    }
  };
 

  useEffect(() => {
    fetchData();
    fetchConn();
   
  }, []);

  const handleInputChange = (e) => {
    setNewItem({ ...newItem, [e.target.name]: e.target.value });
  };

  const handleAddItem = async () => {
    try {
      const token = localStorage.getItem('token');
      

      const response = await fetch(`${API_URL}/metric`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(newItem),
      });

      if (response.ok) {
        // Reset the form and fetch the updated data
        setNewItem({
          project: '',
          catalog: '',
          item_desc: '',
          item_condition: '',
          conn_name:'',
          dashboard_url: '',
        });
        fetchData();
        fetchConn();
      } else {
        console.error('Error adding item:', response.status);
      }
    } catch (error) {
      console.error('Error adding item:', error);
    } finally {
      onClose();
    }
  };

  const handleUpdateItem = async () => {

    console.log(editItem)



    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`${API_URL}/metric/${editItem.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(editItem),
      });

      if (response.ok) {
        // Fetch the updated data
        fetchData();
        fetchConn();
      } else {
        console.error('Error updating item:', response.status);
      }
    } catch (error) {
      console.error('Error updating item:', error);
    } finally {
      setEditItem(null);
      onClose();
    }
  }


  const handleDeleteItem = async (id) => {



    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`${API_URL}/metric/${id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(editItem),
      });

      if (response.ok) {
        // Fetch the updated data
        fetchData();
        fetchConn();
      } else {
        console.error('Error updating item:', response.status);
      }
    } catch (error) {
      console.error('Error updating item:', error);
    } finally {
      setConfirmOpen(false)
    }
  }


  const columns = useMemo(
    () => [
      { Header: 'ID', accessor: 'id' },
      { Header: 'Project', accessor: 'project' },
      { Header: 'Catalog', accessor: 'catalog' },
      { Header: 'Item Description', accessor: 'item_desc' },
      { Header: 'Item Condition', accessor: 'item_condition' },
      { Header: 'Connection', accessor: 'conn_name' },
      {
        Header: 'Dashboard URL',
        accessor: 'dashboard_url',
        Cell: ({ value }) => <Box w="100px"><a href={value} target="_blank" rel="noopener noreferrer">{value}</a></Box>,
      },
      {
        Header: 'Status',
        accessor: 'status',
        Cell: ({ value }) => <Box w="80px"><span>{value ? 'Pass' : 'Not Pass'}</span></Box>,
      },
      { Header: 'Status Description', accessor: 'status_desc', Cell: ({ value }) => <Box w="300px">{value!==""?value:"Not have enough information"}</Box> },
      {
        Header: 'Screen', accessor: 'screen',
        Cell: ({ value }) => (
          <Link isExternal href={`/${value}`} target="_blank" rel="noopener noreferrer">
            {value!==""?<Image src={`/${value}`} alt="Screen" maxW="200px" />:<div>Not have screen image</div>}
          </Link>
        ),
      },
      { Header: 'Check Date', accessor: 'check_date' },
      {
        Header: 'ACTION',
        Cell: ({ row }) => <ActionCell
          row={row}
          handleRunClick={fetchData}
          setEditItem={setEditItem}
          onOpen={onOpen}
          setConfirmOpen={setConfirmOpen}
          setDeleteItem={setDeleteItem}

        />,
      },
    ],
    []
  );


  const handleConnectionChange = (e) => {
    if (editItem){
      setEditItem({...editItem,conn_name:e.target.value})
    }else{
      setNewItem({...newItem,conn_name:e.target.value})
    }
};
 
 

  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable({ columns, data});

  return (
    <>
      <Box mt="500h" maxW="90vw" ml="50px" >
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
      </Box>
      <Flex justifyContent="center" mt={4} mb={4}>
        {isLoading ? (
          <Spinner mr={6} size="md" />
        ) : (
          <Button mr={2} onClick={fetchData}>Refresh</Button>
        )}

        <Button mr={2} onClick={() => {
          setEditItem(null)
          onOpen()
        }
        }>Add Metric</Button>

        
        
      </Flex>
      <AlertDialog isOpen={confirmIsOpen} leastDestructiveRef={cancelRef} onClose={() => { setConfirmOpen(false) }}>
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Item
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to delete this item? <hr />
              {deleteItem?.id},{deleteItem?.project},{deleteItem?.catalog},<br />{deleteItem?.item_desc}<br />{deleteItem?.item_condition}

            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={() => { setConfirmOpen(false) }}>
                Cancel
              </Button>
              <Button colorScheme="red" ml={3} onClick={() => handleDeleteItem(deleteItem?.id)}>
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
              {editItem ? 'Edit Item' : 'Add Item'}
            </AlertDialogHeader>
            <AlertDialogBody>
            Project<Input
                name="project"
                placeholder="Project"
                value={editItem ? editItem.project : newItem.project}
                onChange={(e) => {
                  if (editItem) {
                    setEditItem({ ...editItem, project: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              />
              Connections
              {connections.length>0&&<Select placeholder="Select an option" value={(editItem ? editItem.conn_name : newItem.conn_name)} onChange={handleConnectionChange} >
              {connections.map((option) => (
                <option
                  key={option.conn_name}
                  value={option.conn_name}
                 //selected={option.conn_name === (editItem ? editItem.conn_name : newItem.conn_name)}
                >
                  {option.conn_name}
                </option>
              ))}
            </Select>}
              
              Catalog
              <Input
                name="catalog"
                placeholder="Catalog"
                value={editItem ? editItem.catalog : newItem.catalog}
                onChange={(e) => {
                  if (editItem) {
                    setEditItem({ ...editItem, catalog: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              />Description
              <Input
                name="item_desc"
                placeholder="Item Description"
                value={editItem ? editItem.item_desc : newItem.item_desc}
                onChange={(e) => {
                  if (editItem) {
                    setEditItem({ ...editItem, item_desc: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              /> Condition
              <Input
                name="item_condition"
                placeholder="Item Condition"
                value={editItem ? editItem.item_condition : newItem.item_condition}
                onChange={(e) => {
                  if (editItem) {
                    setEditItem({ ...editItem, item_condition: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              />Dashboard URL
              <Input
                name="dashboard_url"
                placeholder="Dashboard URL"
                value={editItem ? editItem.dashboard_url : newItem.dashboard_url}
                onChange={(e) => {
                  if (editItem) {
                    setEditItem({ ...editItem, dashboard_url: e.target.value });
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
              {editItem ? (
                <Button colorScheme="blue" onClick={handleUpdateItem} ml={3}>
                  Update
                </Button>
              ) : (
                <Button colorScheme="blue" onClick={handleAddItem} ml={3}>
                  Add
                </Button>
              )}
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
};

export default TableComponent;