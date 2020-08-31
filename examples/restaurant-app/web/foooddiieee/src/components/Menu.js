import React, { useState, useContext } from 'react';
import { useQuery, useMutation } from '@apollo/react-hooks';
import gql from "graphql-tag";
import './Foooddiieee.css';
import { Container } from 'react-bootstrap';
import { Row } from 'react-bootstrap';
import { Col } from 'react-bootstrap';
import { Mutation } from 'react-apollo'
import { AppContext } from './AppContext'

const READ_INFO = gql`
  query items {
    items {
      name
      itemType
    }
  }
`;

const CHANGE_MENU = gql`
  mutation mutateItem($itemType: ItemTypes!, $name: String!, $action: Actions!){
    mutateItem(itemType: $itemType, name: $name, action: $action) {
      itemType name action
    }
  }
`;

function Menu() {
  let itemType, name;
  const [isManager, setIsManager] = useContext(AppContext);
  const { data, loading, error } = useQuery(READ_INFO);
  const [mutateItem] = useMutation(CHANGE_MENU, {
    refetchQueries: ["items"]
  });

  if (loading) return <p>loading...</p>;
  if (error) return <p>ERROR</p>;
  if (!data) return <p>Not found</p>;

  if(isManager) {
    return (
     <section className="about-section">
       <Container>
         <Row>
            <Col className="centered">
              <h1>Menu</h1>
             </Col>
         </Row>
          <Row>
            <Col className="opening-time">
              <h3>
                Appetizer
              </h3>
              <ul>
                {data.items.filter(item => item.itemType === 'Appetizer').map((fItem) =>
                  <li key={fItem.name} className="w-100">
                    <Row>
                      <Col className="opening-time">        
                        <span>{fItem.name}</span>
                      </Col>
                      <Col className="opening-time">
                        <button className="btn btn-sm btn-danger rounded-circle" onClick={() => {
                          mutateItem({ variables: { name: fItem.name, itemType: fItem.itemType,  action: 'Remove'} });
                        }}>X</button>  
                      </Col>
                    </Row>                
                  </li>
                )}
              </ul>
            </Col>
            <Col className="opening-time">
              <h3>
                Entree
              </h3>
              <ul>
                {data.items.filter(item => item.itemType === 'Entree').map((fItem) =>
                  <li key={fItem.name} className="w-100">
                    <Row>
                      <Col className="opening-time">        
                        <span>{fItem.name}</span>
                      </Col>
                      <Col className="opening-time">
                        <button className="btn btn-sm btn-danger rounded-circle" onClick={() => {
                          mutateItem({ variables: { name: fItem.name, itemType: fItem.itemType,  action: 'Remove'} });
                        }}>X</button>  
                      </Col>
                    </Row>                   
                  </li>
                )}
              </ul>
            </Col>
            <Col className="opening-time">
              <h3>
                Dessert
              </h3>
              <ul>
                {data.items.filter(item => item.itemType === 'Dessert').map((fItem) =>
                  <li key={fItem.name} className="w-100">
                    <Row>
                      <Col className="opening-time">        
                        <span>{fItem.name}</span>
                      </Col>
                      <Col className="opening-time">
                        <button className="btn btn-sm btn-danger rounded-circle" onClick={() => {
                          mutateItem({ variables: { name: fItem.name, itemType: fItem.itemType,  action: 'Remove'} });
                        }}>X</button>  
                      </Col>
                    </Row>                    
                  </li>
                )}
              </ul>
            </Col>
          </Row>
         </Container>
        <form onSubmit={e => {
          e.preventDefault();
          mutateItem({ variables: { name: name.value,  itemType: itemType.value, action: 'Add'} });
          itemType.value = '';
          name.value = '';
        }}>
          <div class="input-group">
            <div class="input-group-prepend">
              <span class="input-group-text" id="">Add a new menu item</span>
            </div>
          <input className="form-control" type="text" placeholder="Enter item type" ref={node => { itemType = node; }}></input>
          <input className="form-control" type="text" placeholder="Enter name" ref={node => { name = node; }}></input>
          <button className="btn btn-outline-secondary" type="submit">Submit</button>
          </div>
        </form>       
       </section>
    );    
  } else {
    return (
     <section className="about-section">
       <Container>
         <Row>
            <Col className="centered">
              <h1>Menu</h1>
             </Col>
         </Row>
          <Row>
            <Col className="opening-time">
              <h3>
                Appetizer
              </h3>
              <ul>
                {data.items.filter(item => item.itemType === 'Appetizer').map((fItem) =>
                  <li key={fItem.name} className="w-100">
                    <span>{fItem.name}</span>                 
                  </li>
                )}
              </ul>
            </Col>
            <Col className="opening-time">
              <h3>
                Entree
              </h3>
              <ul>
                {data.items.filter(item => item.itemType === 'Entree').map((fItem) =>
                  <li key={fItem.name} className="w-100">
                    <span>{fItem.name}</span>                  
                  </li>
                )}
              </ul>
            </Col>
            <Col className="opening-time">
              <h3>
                Dessert
              </h3>
              <ul>
                {data.items.filter(item => item.itemType === 'Dessert').map((fItem) =>
                  <li key={fItem.name} className="w-100">
                    <span>{fItem.name}</span>                  
                  </li>
                )}
              </ul>
            </Col>
          </Row>
         </Container>      
       </section>
    );     
  }

}

export default Menu;