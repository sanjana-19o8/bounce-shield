import React from 'react';

const Card = ({ children, className }) => {
  return (
    <div className={`rounded-2xl shadow-md border border-gray-200 ${className}`}>
      {children}
    </div>
  );
};

export const CardContent = ({ children }) => {
  return <div className="p-4">{children}</div>;
};

export default Card;
