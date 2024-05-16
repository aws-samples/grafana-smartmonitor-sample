import { Spinner, Button, Box } from "@chakra-ui/react";

import { useState } from "react";

import { API_URL } from '../constants';

const ActionCell = ({ row, handleRunClick,setEditItem,onOpen,setDeleteItem,setConfirmOpen}) => {
    const [isRunning, setIsRunning] = useState(false);

    const handleClick = async (rowId) => {
      setIsRunning(true);

      try {
        const token = localStorage.getItem('token');
        const headers = {
          'Authorization': `Bearer ${token}`,
        };
        const response = await fetch(`${API_URL}/run/${rowId}`, {
          method: 'POST',
          headers
        });
        if (response.ok) {
          // Refresh items after successful API call
          handleRunClick();
        } else {
          console.error('Error running API:', response.status);
        }
      } catch (error) {
        console.error('Error running API:', error);
      } finally {
        setIsRunning(false)
      }

    };

    const handleEdit = (item) => {
      setEditItem(item);
      onOpen();
    };


    const handleDelete = (item) => {
      setDeleteItem(item);
      setConfirmOpen(true);
    };

    return (
      <Box>
        {isRunning ? (
          <Spinner mr={6} size="md" />
        ) : (
          <Button mb={"2"} colorScheme="blue" onClick={() => { handleClick(row.original.id) }}>
            Run
          </Button>
        )}
        <Button mb={"2"} colorScheme="green" onClick={() => handleEdit(row.original)}>
          Edit
        </Button>
        {/* <Button mb={"2"} colorScheme="orange" onClick={() => {console.log("history ")}}>
          History
        </Button> */}
        <Button colorScheme="red" onClick={() => { handleDelete(row.original) }}>Delete</Button>
      </Box>
    );
  };


  
export default ActionCell

