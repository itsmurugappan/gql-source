import React, { useState, createContext } from 'react';

export const AppContext = createContext();

export const AppContextProvider = props => {
  const [isManager, setIsManager] = useState(false);
  return (
    <AppContext.Provider value={[isManager, setIsManager]}>
      {props.children}
    </AppContext.Provider>
  );
}