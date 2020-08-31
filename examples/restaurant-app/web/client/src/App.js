import React from 'react';
import './App.css';
import { useQuery, useMutation } from '@apollo/react-hooks';
import gql from "graphql-tag";

const CREATE_MENU = gql`
  mutation {
    mutateItem(itemType: Appetizer, name: "fries", action: Add) {
      name
    }
  }
`;

const REMOVE_MENU = gql`
  mutation {
    mutateItem(itemType: Appetizer, name: "fries", action: Remove) {
      name
    }
  }
`;

const UPDATE_INFO = gql`
  mutation {
    mutateItem(itemType: Appetizer, name: "fries", action: Remove) {
      name
    }
  }
`;

// const UPDATE_TODO = gql`
//   mutation UpdateTodo($id: String!) {
//     updateTodo(id: $id)
//   }
// `;

function App() {
  let input;
  const { data, loading, error } = useQuery(READ_INFO);
  const [createMenu] = useMutation(CREATE_MENU);
  const [deleteMenu] = useMutation(REMOVE_MENU);

  if (loading) return <p>loading...</p>;
  if (error) return <p>ERROR</p>;
  if (!data) return <p>Not found</p>;

  return (
    <div className="app">
      <h3>New Menu Item</h3>
      <form onSubmit={e => {
        e.preventDefault();
        createMenu({ variables: { text: input.value } });
        input.value = '';
        window.location.reload();
      }}>
        <input className="form-control" type="text" placeholder="Enter todo" ref={node => { input = node; }}></input>
        <button className="btn btn-primary px-5 my-2" type="submit">Submit</button>
      </form>
      <ul>
        {data.items.map((item) =>
          <li key={item.name} className="w-100">
            <span>{item.name}</span>
          </li>
        )}
      </ul>
    </div>
  );
}

export default App;