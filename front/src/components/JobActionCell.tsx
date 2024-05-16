import { Spinner, Button, Box } from "@chakra-ui/react";

import { useState } from "react";

const JobActionCell = ({ row, handleRunClick,setEditItem,onOpen,setDeleteItem,setConfirmOpen}) => {
    

    const handleEdit = (item) => {
      console.log(item)
      setEditItem(item);
      onOpen();
    };


    const handleDelete = (item) => {
      console.log(item)
      setDeleteItem(item);
      setConfirmOpen(true);
    };

    return (
      <Box>
        
        
        <Button ml={"2"} colorScheme="green" onClick={() => handleEdit(row.original)}>
          Edit
        </Button>
        
        <Button ml={"2"} colorScheme="red" onClick={() => { handleDelete(row.original) }}>Delete</Button>
      </Box>
    );
  };


  export default JobActionCell