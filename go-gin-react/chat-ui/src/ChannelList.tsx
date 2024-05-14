
import { useState, useEffect } from 'react';
import { useParams } from "react-router-dom";

interface SelectedChannel {
  name: string;
  id: number;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const ChannelsList = ({ setSelectedChannel }: { setSelectedChannel: (channel: SelectedChannel) => void }) => {
  const { channelId } = useParams<{ channelId: string }>(); // Adjust type accordingly
  const [channels, setChannels] = useState<SelectedChannel[]>([]); // Set initial type to Channel[]
  const [newChannelName, setNewChannelName] = useState('');

  useEffect(() => {
    if (channelId) {
      const channel = channels.find((channel) => channel.id === parseInt(channelId));
      if (channel) {
        setSelectedChannel({ name: channel.name, id: parseInt(channelId) });
      }
    }
  }, [channelId, channels, setSelectedChannel]);

  useEffect(() => {
    const fetchChannels = async () => {
      const response = await fetch('http://localhost:8080/channels');
      const data: SelectedChannel[] = await response.json(); // Specify the return type
      setChannels(data || []);
    };
    fetchChannels();
  }, []);

  const handleAddChannel = async () => {
    const response = await fetch('http://localhost:8080/channels', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: newChannelName }),
    });

    if (response.ok) {
      const newChannel: SelectedChannel = await response.json(); // Specify the return type
      setChannels([...channels, { id: newChannel.id, name: newChannelName }]);
      setNewChannelName('');
    }
  };

  return (
    <div className="flex flex-col h-full bg-gray-100 border-r">
      <div className="bg-gray-700 text-white p-2">
        Channels
      </div>
      <div className="overflow-y-auto flex-grow p-4">
        {
          // channels ? (
          //   <ul className="w-full">
          //     {channels.map((channel) => (
          //       <li
          //         key={channel.id}
          //         className={`p-2 rounded-md w-full cursor-pointer ${parseInt(channelId) === channel.id ? 'bg-blue-500 text-white' : 'hover:bg-gray-200'}`}
          //         onClick={() => setSelectedChannel(channel)}
          //       >
          //         {channel.name}
          //       </li>
          //     ))}
          //   </ul>
          channels ? (
            <ul className="w-full">
              {channels.map((channel) => (
                <li
                  key={channel.id}
                  className={`p-2 rounded-md w-full cursor-pointer ${channelId !== undefined && parseInt(channelId) === channel.id
                    ? 'bg-blue-500 text-white'
                    : 'hover:bg-gray-200'
                    }`}
                  onClick={() => setSelectedChannel(channel)}
                >
                  {channel.name}
                </li>
              ))}
            </ul>
          ) : (
            <div className="text-center text-gray-600">
              Please add a Channel
            </div>
          )}
      </div>
      <div className="flex flex-col p-4">
        <input
          type="text"
          value={newChannelName}
          onChange={(e) => setNewChannelName(e.target.value)}
          placeholder="New channel..."
          className="mb-4 p-2 w-full border rounded-md bg-white"
        />
        <button onClick={handleAddChannel} className="p-2 bg-blue-500 text-white rounded-md">Add Channel</button>
      </div>
    </div>
  );
};

export default ChannelsList;
