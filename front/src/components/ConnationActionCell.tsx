import { Spinner, Button, Box } from "@chakra-ui/react";

import { useState } from "react";

const ConnationActionCell = ({ row, setEditItem,onOpen,setDeleteItem,setConfirmOpen}) => {
  
    const handleEdit = (item) => {
      setEditItem(item);
      onOpen();
    };
    const handleDelete = (item) => {
      setDeleteItem(item);
      setConfirmOpen(true);
    };

    return (
      <Box w={180}>
        <Button mt={"2"} ml={"2"} colorScheme="green" onClick={() => handleEdit(row.original)}>
          Edit
        </Button>
        <Button mt={"2"} ml={"2"} colorScheme="red" onClick={() => { handleDelete(row.original) }}>Delete</Button>
      </Box>
    );
  };


  export default ConnationActionCell