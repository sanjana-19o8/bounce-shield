import React, { useState } from 'react';
import axios from 'axios';
import { Card, CardContent } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { toast } from 'react-hot-toast';
import Sidebar from '../components/Sidebar';

const TempEmailDashboard = () => {
  const [tempEmail, setTempEmail] = useState('');
  const [inbox, setInbox] = useState([]);
  const [loading, setLoading] = useState(false);

  const createTempEmail = async () => {
    setLoading(true);
    try {
      const response = await axios.get('http://localhost:8082/api/create-temp-email');
      setTempEmail(response.data.Address);
      toast.success('Temporary email created!');
    } catch (error) {
      toast.error('Failed to create temp email');
    }
    setLoading(false);
  };

  const fetchInbox = async () => {
    if (!tempEmail) {
      toast.error('Create a temp email first!');
      return;
    }
    setLoading(true);
    try {
      const response = await axios.get(`http://localhost:8082/api/check-temp-inbox?tempEmail=${tempEmail}`);
      setInbox(response.data);
      toast.success('Inbox updated');
    } catch (error) {
      toast.error('Failed to fetch inbox');
    }
    setLoading(false);
  };

  return (
    <div className="flex">
      <Sidebar />
      <div className="p-6 w-full">
        <h1 className="text-2xl font-semibold mb-4">📬 Temporary Email Dashboard</h1>
        <div className="space-y-4">
          <Button onClick={createTempEmail} disabled={loading} className="bg-blue-600 hover:bg-blue-700">Generate Temp Email</Button>
          {tempEmail && (
            <div>
              <p className="text-lg font-medium">Your Temp Email:</p>
              <Card className="bg-gray-100 mt-2 p-4">
                <CardContent>{tempEmail}</CardContent>
              </Card>
              <Button onClick={fetchInbox} className="mt-4 bg-green-600 hover:bg-green-700">Check Inbox</Button>
            </div>
          )}

          {inbox.length > 0 && (
            <div className="mt-6">
              <h2 className="text-xl font-semibold mb-2">📥 Inbox Messages</h2>
              <div className="space-y-2">
                {inbox.map((msg, idx) => (
                  <Card key={idx} className="bg-white p-3">
                    <CardContent>{msg}</CardContent>
                  </Card>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default TempEmailDashboard;
