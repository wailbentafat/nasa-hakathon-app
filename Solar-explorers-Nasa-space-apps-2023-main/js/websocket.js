
  const socket = new WebSocket('ws://localhost:8080/ws');

  
  socket.addEventListener('open', (event) => {
      console.log('WebSocket connection opened');
      document.getElementById('resultText').innerText = 'WebSocket connection opened.';
  });

 
  socket.addEventListener('move ', (event) => {
      console.log('Message from server: ', event.data);
      document.getElementById('resultText').innerText = 'Server says: ' + event.data;
  });

 
  socket.addEventListener('close', (event) => {
      console.log('WebSocket connection closed');
      document.getElementById('resultText').innerText = 'WebSocket connection closed.';
  });

 
  socket.addEventListener('error', (event) => {
      console.error('WebSocket error observed:', event);
      document.getElementById('resultText').innerText = 'WebSocket error occurred.';
  });

 
  