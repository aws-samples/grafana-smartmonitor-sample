'use client'
import React, { useContext, useEffect, useMemo, useState } from 'react';
import {  Table, Thead, Tbody, Tr, Th, Td, Input, Button, Box, Flex, Spinner, AlertDialog, AlertDialogBody, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay, useDisclosure, Checkbox } from '@chakra-ui/react';
import { useTable } from 'react-table';

import JobActionCell from "@components/JobActionCell"
import JobStatusCell from '@components/JobStatusCell';
import { AuthContext } from '@components/AuthContext';

import { useRouter } from 'next/navigation';

import { API_URL } from '../../constants';


const SchedulerPage = () => {
  const [data, setData] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  const { isOpen, onOpen, onClose } = useDisclosure();

  const [confirmIsOpen, setConfirmOpen] = useState(false);

  const cancelRef = React.useRef();

  const { isAuthenticated } = useContext(AuthContext);

  const route = useRouter()

  const [newJob, setNewJob] = useState({
    project: '',
    cron: '',
    enable: false,
  });

  const [editJob, setEditJob] = useState(null);
  const [deleteJob, setDeleteJob] = useState(null);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const token = localStorage.getItem('token');
      // Construct the headers object with the Authorization header
      const headers = {
        'Authorization': `Bearer ${token}`,
        // Add any other headers if needed
      };
      const response = await fetch(`${API_URL}/jobs`,{headers});
      const data = await response.json();
      if (data.jobs) {
        setData(data.jobs);
      }

    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setTimeout(() => setIsLoading(false), 1000);
    }
  };

  
  useEffect(() => {
    fetchData();
    
  }, []);

  const handleInputChange = (e) => {
    setNewJob({ ...newJob, [e.target.name]: e.target.value });
  };

  const handleAddJob = async () => {
    try {
      const token = localStorage.getItem('token');
      
      const response = await fetch(`${API_URL}/job`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(newJob),
      });

      if (response.ok) {
        // Reset the form and fetch the updated data
        setNewJob({
          project: '',
          cron: '',
          enable: false,
        });
        fetchData();
      } else {
        console.error('Error adding job:', response.status);
      }
    } catch (error) {
      console.error('Error adding job:', error);
    } finally {
      onClose();
    }
  };

  const handleUpdateJob = async () => {
    console.log(editJob)

    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`${API_URL}/job/${editJob.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(editJob),
      });

      if (response.ok) {
        // Fetch the updated data
        fetchData();
      } else {
        console.error('Error updating job:', response.status);
      }
    } catch (error) {
      console.error('Error updating job:', error);
    } finally {
      setEditJob(null);
      onClose();
    }
  }

  const handleDeleteJob = async (id) => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`${API_URL}/job/${id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        // Fetch the updated data
        fetchData();
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
      { Header: 'Project', accessor: 'project' },
      { Header: 'Cron', accessor: 'cron' },
      { Header: 'Enable', accessor: 'enable', Cell: ({ value }) => <Box w="80px"><span>{value ? 'Enabled' : 'Disabled'}</span></Box> },
      { Header: 'Status', 
      Cell: ({ row }) => < JobStatusCell
      row={row}
    />

      },
      
      {
        Header: 'ACTION',
        Cell: ({ row }) => <JobActionCell
          row={row}
          handleRunClick={fetchData}
          setEditItem={setEditJob}
          onOpen={onOpen}
          setConfirmOpen={setConfirmOpen}
          setDeleteItem={setDeleteJob}
        />,
      },

      

    ],
    []
  );

  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable({ columns, data });

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
          setEditJob(null)
          onOpen()
        }
        }>Add Job</Button>
        
      </Flex>
      <AlertDialog isOpen={confirmIsOpen} leastDestructiveRef={cancelRef} onClose={() => { setConfirmOpen(false) }}>
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Job
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to delete this job? <hr />
              {deleteJob?.id}, {deleteJob?.project}, {deleteJob?.cron}, {deleteJob?.enable ? 'Enabled' : 'Disabled'}
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={() => { setConfirmOpen(false) }}>
                Cancel
              </Button>
              <Button colorScheme="red" ml={3} onClick={() => handleDeleteJob(deleteJob?.id)}>
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
              {editJob ? 'Edit Job' : 'Add Job'}
            </AlertDialogHeader>
            <AlertDialogBody>
              Project<Input
                name="project"
                placeholder="Project"
                value={editJob ? editJob.project : newJob.project}
                onChange={(e) => {
                  if (editJob) {
                    setEditJob({ ...editJob, project: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              />Cron

              
              <Input
                name="cron"
                placeholder="Cron"
                value={editJob ? editJob.cron : newJob.cron}
                onChange={(e) => {
                  if (editJob) {
                    setEditJob({ ...editJob, cron: e.target.value });
                  } else {
                    handleInputChange(e);
                  }
                }}
                mb={2}
              />Enable
              <br/>
              <Checkbox
                name="enable"
                isChecked={editJob ? editJob.enable : newJob.enable}
                onChange={(e) => {
                  if (editJob) {
                    setEditJob({ ...editJob, enable: e.target.checked });
                  } else {
                    setNewJob({ ...newJob, enable: e.target.checked });
                  }
                }}
                mb={2}
              >
                Enabled
              </Checkbox>
              
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              {editJob ? (
                <Button colorScheme="blue" onClick={handleUpdateJob} ml={3}>
                  Update
                </Button>
              ) : (
                <Button colorScheme="blue" onClick={handleAddJob} ml={3}>
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

export default SchedulerPage;
