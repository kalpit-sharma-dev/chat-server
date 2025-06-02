// BankingChatbotUI.tsx

import React, { useState, useRef, useEffect } from 'react';
import {
  View,
  TextInput,
  Button,
  FlatList,
  Text,
  StyleSheet,
  Image,
  ActivityIndicator,
  TouchableOpacity
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';
import * as Speech from 'expo-speech';
import * as SpeechToText from 'expo-speech';
import Markdown from 'react-native-markdown-display';

interface Message {
  type: 'user' | 'bot';
  text: string;
  timestamp: string;
}

const ChatScreen = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const flatListRef = useRef<FlatList>(null);
  const controllerRef = useRef<AbortController | null>(null);

  const getTime = () => new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

  const scrollToEnd = () => {
    setTimeout(() => flatListRef.current?.scrollToEnd({ animated: true }), 100);
  };

  useEffect(() => {
    loadMessages();
  }, []);

  useEffect(() => {
    saveMessages();
    scrollToEnd();
  }, [messages]);

  const saveMessages = async () => {
    await AsyncStorage.setItem('chatMessages', JSON.stringify(messages));
  };

  const loadMessages = async () => {
    const stored = await AsyncStorage.getItem('chatMessages');
    if (stored) setMessages(JSON.parse(stored));
  };

  const parseSSEChunk = (chunk: string) => {
    return chunk.split('\n').filter(line => line.startsWith('data:')).map(line => line.replace(/^data: /, '')).join('');
  };

  const handleSend = async () => {
    if (!input.trim()) return;

    const userMsg: Message = { type: 'user', text: input, timestamp: getTime() };
    const botMsg: Message = { type: 'bot', text: '', timestamp: getTime() };

    const newMessages = [...messages, userMsg, botMsg];
    setMessages(newMessages);
    setInput('');
    setIsTyping(true);

    const messageIndex = newMessages.length - 1;

    controllerRef.current = new AbortController();

    const response = await fetch('https://your-api-endpoint.com/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'text/event-stream',
      },
      body: JSON.stringify({ message: input }),
      signal: controllerRef.current.signal,
    });

    const reader = response.body?.getReader();
    const decoder = new TextDecoder('utf-8');

    if (!reader) return;

    let botMessage = '';
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      const chunk = decoder.decode(value, { stream: true });
      const parsedChunk = parseSSEChunk(chunk);
      botMessage += parsedChunk;

      setMessages(prev =>
        prev.map((msg, idx) =>
          idx === messageIndex ? { ...msg, text: botMessage } : msg
        )
      );
    }

    setIsTyping(false);
    Speech.speak(botMessage);
  };

  const handleVoiceInput = () => {
    SpeechToText.stop(); // placeholder - replace with real voice-to-text API
    alert('Voice input feature requires speech-to-text SDK integration');
  };

  const renderMessage = ({ item }: { item: Message }) => (
    <View
      style={[
        styles.messageContainer,
        item.type === 'user' ? styles.userContainer : styles.botContainer,
      ]}
    >
      {item.type === 'bot' ? (
        <Image source={require('./bot-avatar.png')} style={styles.avatar} />
      ) : (
        <View style={styles.userAvatar}><Text style={styles.userInitials}>U</Text></View>
      )}
      <View style={styles.messageBubble}>
        <Markdown style={markdownStyles}>{item.text}</Markdown>
        <Text style={styles.timestamp}>{item.timestamp}</Text>
      </View>
    </View>
  );

  return (
    <View style={styles.container}>
      <FlatList
        ref={flatListRef}
        data={messages}
        keyExtractor={(_, index) => index.toString()}
        renderItem={renderMessage}
        contentContainerStyle={{ padding: 10 }}
      />
      {isTyping && (
        <View style={styles.typingIndicator}>
          <ActivityIndicator size="small" color="#555" />
          <Text style={{ marginLeft: 8 }}>BankBot is typing...</Text>
        </View>
      )}
      <View style={styles.inputContainer}>
        <TouchableOpacity onPress={handleVoiceInput} style={styles.micButton}>
          <Text style={{ fontSize: 18 }}>🎤</Text>
        </TouchableOpacity>
        <TextInput
          style={styles.input}
          placeholder="Type your message..."
          value={input}
          onChangeText={setInput}
          onSubmitEditing={handleSend}
        />
        <Button title="Send" onPress={handleSend} />
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f9f9f9' },
  inputContainer: {
    flexDirection: 'row',
    padding: 10,
    borderTopWidth: 1,
    borderColor: '#ddd',
    backgroundColor: '#fff',
    alignItems: 'center'
  },
  micButton: {
    padding: 10
  },
  input: {
    flex: 1,
    borderColor: '#ccc',
    borderWidth: 1,
    borderRadius: 8,
    padding: 10,
    marginHorizontal: 10,
    backgroundColor: '#fff',
  },
  messageContainer: {
    flexDirection: 'row',
    marginVertical: 4,
    alignItems: 'flex-end',
  },
  userContainer: {
    justifyContent: 'flex-end',
    alignSelf: 'flex-end',
  },
  botContainer: {
    alignSelf: 'flex-start',
  },
  avatar: {
    width: 30,
    height: 30,
    borderRadius: 15,
    marginRight: 8,
  },
  userAvatar: {
    width: 30,
    height: 30,
    borderRadius: 15,
    backgroundColor: '#87ceeb',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 8,
  },
  userInitials: {
    color: '#fff',
    fontWeight: 'bold',
  },
  messageBubble: {
    maxWidth: '75%',
    backgroundColor: '#ECECEC',
    borderRadius: 10,
    padding: 10,
  },
  timestamp: {
    fontSize: 10,
    color: '#666',
    marginTop: 4,
    alignSelf: 'flex-end',
  },
  typingIndicator: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 10,
    marginLeft: 10,
  },
});

const markdownStyles = {
  text: { fontSize: 16 },
};

export default ChatScreen;
