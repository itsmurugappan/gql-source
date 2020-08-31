import React, { useState } from 'react';
import { useQuery, useMutation } from '@apollo/react-hooks';
import gql from "graphql-tag";
import './Foooddiieee.css';


const UPDATE_INFO = gql`
  mutation UpdateInfo($hours: String!, $address: String!){
    updateInfo(hours: $hours, address: $address) {
      hours address
    }
  }
`;

function Manage() {
  const [hours, setHours] = useState('')
  const [address, setAddress] = useState('')
  const [updateInfo] = useMutation(UPDATE_INFO, {
    variables: {hours, address}
  }); 

  return (
    <div className="app">
      <h3>Update Info</h3>
      <input 
        onChange={e => setHours(e.target.value)}
      />
      <input 
        onChange={e => setAddress(e.target.value)}
      />
      <button onClick={updateInfo}>Update Info</button>     
    </div>
  );
}

export default Manage;