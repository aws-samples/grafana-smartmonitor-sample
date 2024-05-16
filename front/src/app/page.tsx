'use client';

import { Box ,Text} from '@chakra-ui/react';

const Home = () => {
  return  <Box mt="500h" maxW="800px" m="0 auto" p={5}>
    <Text fontSize="2xl" fontFamily="monospace" fontWeight="bold">
          Grafana SmartMonitor
    </Text>
    <br></br>
    <Text fontSize="xl" fontFamily="monospace" >
        An AI-powered monitoring solution that intelligently detects and alerts anomalies, predicts failures, and provides actionable insights for proactive maintenance and optimal system performance. <br/>
        
    </Text>
    <hr/><br/>
    Powered By Amazon Bedrock and Anthropic Claude3
    </Box>;
};

export default Home;
