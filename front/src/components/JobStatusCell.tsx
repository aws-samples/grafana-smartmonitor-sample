import { Button, Box, Modal,  ModalBody, ModalCloseButton, ModalContent, ModalFooter, ModalHeader, ModalOverlay, Textarea } from "@chakra-ui/react";

import { useContext, useEffect, useState } from "react";

import { API_URL } from '../constants';
import { AuthContext } from "./AuthContext";


const JobStatusCell = ({ row}) => {
    
  const [summaryData, setSummaryData] = useState(null);
  const [isOpen, setIsOpen] = useState(false)
  

  useEffect(() => {
    const fetchData = async () => { 
      
      try {
        if (row.original.enable === true) {
          const token = localStorage.getItem('token');
          // Construct the headers object with the Authorization header
          const headers = {
            'Authorization': `Bearer ${token}`,
          };

          const response = await fetch(`${API_URL}/status/${row.original.project}`,{headers});
          const data = await response.json();
          const summary = JSON.parse(data.job.Summary);
          setSummaryData({
            health: summary.health,
            summary: summary.summary,
            metrics:summary.metrics,
          });
        }
      } catch (error) {
        console.error("Error fetching or parsing data:", error);
        // You can also display an error message to the user if needed
      }
      
    };

    const interval = setInterval(fetchData, 10000); // refresh data every 1 min

    return () => clearInterval(interval); 
  }, [row.original.project,row.original.enable]); 
  const handleResultClick = () => {
    if (summaryData) {
      setIsOpen(true)
    }
  };

 
  const getHealthColor = (health) => {
    switch (health) {
      case "Good":
      case "Very Good":
        return "green";
      case "Bad":
        return "red";
      case "Very Bad":
          return "red";
      default:
        return "gray";
    }
  };

    return (
      <Box>
        
        {summaryData && (
        <Button ml={"2"} colorScheme={getHealthColor(summaryData.health)} onClick={handleResultClick}>
          {summaryData.health}
        </Button>
      )}
       
        <Modal  isOpen={isOpen} onClose={()=>{setIsOpen(false)}}>
        <ModalOverlay />
        <ModalContent >
          
          <ModalHeader>Project Summary</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            {summaryData && (
              <>
                <Textarea value={JSON.stringify(summaryData, null, 2)} isReadOnly  height="300px"/>
              </>
            )}
          </ModalBody>

          <ModalFooter>
            <Button colorScheme="blue" mr={3} onClick={()=>setIsOpen(false)}>
              Close
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
      </Box>
    );
  };


  export default JobStatusCell