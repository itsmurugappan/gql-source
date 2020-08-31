import React, { useContext } from 'react';
import { useQuery, useMutation } from '@apollo/react-hooks';
import gql from "graphql-tag";
import './Foooddiieee.css';
import { Container } from 'react-bootstrap';
import { Row } from 'react-bootstrap';
import { Col } from 'react-bootstrap';
import { AppContext } from './AppContext'

const READ_INFO = gql`
  query allInfo {
    allInfo {
      address
      hours
    }
  }
`;

const UPDATE_INFO = gql`
  mutation UpdateInfo($hours: String!, $address: String!){
    updateInfo(hours: $hours, address: $address) {
      hours address
    }
  }
`;

function Info() {
  let hours, address;
  const [isManager, setIsManager] = useContext(AppContext);  
  const { loading, data } = useQuery(READ_INFO);
  const [updateInfo] = useMutation(UPDATE_INFO, {
    variables: {hours, address}, refetchQueries: ["allInfo"]
  }); 

  if (loading) return <p>loading...</p>;
  if (!data) return <p>Not found</p>;

  if(isManager) {
    return (
     <Container>
        <Row>
          <Col className="opening-time">
            <h3 className="info-centered">
              Hours
            </h3>
            <span className="info-centered">{data.allInfo.hours}</span>
          </Col>
          <Col className="opening-time">
            <h3 className="info-centered">
              Address
            </h3>
            <span className="info-centered">{data.allInfo.address}</span>
          </Col>
        </Row>
        <br/ >   
        <form onSubmit={e => {
          e.preventDefault();
          updateInfo({ variables: { hours: hours.value,  address: address.value} });
          hours.value = '';
          address.value = '';
        }}>
          <div class="input-group">
            <div class="input-group-prepend">
              <span class="input-group-text" id="">Update info</span>
            </div>
          <input className="form-control" type="text" placeholder="New hours" ref={node => { hours = node; }}></input>
          <input className="form-control" type="text" placeholder="New address" ref={node => { address = node; }}></input>
          <button className="btn btn-outline-secondary" type="submit">Submit</button>
          </div>
        </form> 
        <br/ >     
       </Container>
    );
  } else {
    return (
     <Container>
        <Row>
          <Col className="opening-time">
            <h3 className="info-centered">
              Hours
            </h3>
            <span className="info-centered">{data.allInfo.hours}</span>
          </Col>
          <Col className="opening-time">
            <h3 className="info-centered">
              Address
            </h3>
            <span className="info-centered">{data.allInfo.address}</span>
          </Col>
        </Row>     
       </Container>
    );    
  }

}

export default Info;