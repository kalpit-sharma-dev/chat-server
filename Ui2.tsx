// App.tsx

import React, { useEffect, useState, useRef } from 'react';
import {
  View,
  TextInput,
  FlatList,
  Text,
  TouchableOpacity,
  Image,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  ScrollView
} from 'react-native';

import AsyncStorage from '@react-native-async-storage/async-storage';
import Voice from '@react-native-voice/voice';
import Tts from 'react-native-tts';
import Markdown from 'react-native-markdown-display';

interface Message {
  id: string;
  text: string;
  sender: 'user' | 'bot';
}

const App = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [streamingText, setStreamingText] = useState('');
  const flatListRef = useRef<FlatList>(null);

  // Load saved messages
  useEffect(() => {
    const loadChat = async () => {
      const history = await AsyncStorage.getItem('chat-history');
      if (history) setMessages(JSON.parse(history));
    };
    loadChat();
  }, []);

  // Save messages
  useEffect(() => {
    AsyncStorage.setItem('chat-history', JSON.stringify(messages));
  }, [messages]);

  // Initialize Voice
  useEffect(() => {
    Voice.onSpeechResults = onSpeechResults;
    return () => Voice.destroy().then(Voice.removeAllListeners);
  }, []);

  const onSpeechResults = (e: any) => {
    const spoken = e.value[0];
    setInput(spoken);
  };

  const handleVoice = async () => {
    try {
      await Voice.start('en-US');
    } catch (e) {
      console.error('Voice error', e);
    }
  };

  const sendMessage = async () => {
    if (!input.trim()) return;

    const userMsg: Message = {
      id: Date.now().toString(),
      text: input,
      sender: 'user',
    };

    setMessages(prev => [...prev, userMsg]);
    setInput('');
    setStreamingText('');

    const botId = Date.now().toString() + '-bot';
    let botMsg: Message = { id: botId, text: '', sender: 'bot' };
    setMessages(prev => [...prev, botMsg]);

    try {
      const response = await fetch('https://your-backend-api.com/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ prompt: userMsg.text }),
      });

      const reader = response.body?.getReader();
      const decoder = new TextDecoder('utf-8');
      let botText = '';

      while (true) {
        const { done, value } = await reader!.read();
        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        botText += chunk;
        setStreamingText(botText);

        setMessages(prev =>
          prev.map(msg => (msg.id === botId ? { ...msg, text: botText } : msg))
        );
      }

      Tts.speak(botText);
    } catch (err) {
      console.error('Error:', err);
    }
  };

  const renderItem = ({ item }: { item: Message }) => (
    <View
      style={[
        styles.message,
        item.sender === 'user' ? styles.userMsg : styles.botMsg,
      ]}
    >
      <View style={styles.avatarContainer}>
        {item.sender === 'user' ? (
          <View style={styles.userAvatar}><Text style={styles.avatarText}>U</Text></View>
        ) : (
          <Image source={require('./assets/bot-avatar.png')} style={styles.botAvatar} />
        )}
      </View>
      <Markdown style={styles.markdown}>{item.text}</Markdown>
    </View>
  );

  return (
    <KeyboardAvoidingView style={styles.container} behavior={Platform.select({ ios: 'padding' })}>
      <FlatList
        ref={flatListRef}
        data={messages}
        keyExtractor={item => item.id}
        renderItem={renderItem}
        onContentSizeChange={() => flatListRef.current?.scrollToEnd({ animated: true })}
        onLayout={() => flatListRef.current?.scrollToEnd({ animated: true })}
      />
      <View style={styles.inputContainer}>
        <TouchableOpacity onPress={handleVoice}>
          <Text style={styles.micButton}>🎤</Text>
        </TouchableOpacity>
        <TextInput
          value={input}
          onChangeText={setInput}
          style={styles.input}
          placeholder="Type your message"
          multiline
        />
        <TouchableOpacity onPress={sendMessage} style={styles.sendButton}>
          <Text style={styles.sendText}>➤</Text>
        </TouchableOpacity>
      </View>
    </KeyboardAvoidingView>
  );
};

export default App;

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f0f4f8' },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    padding: 8,
    backgroundColor: '#fff',
    borderTopWidth: 1,
    borderTopColor: '#ccc',
  },
  input: {
    flex: 1,
    minHeight: 40,
    maxHeight: 100,
    padding: 8,
    backgroundColor: '#eee',
    borderRadius: 20,
    marginHorizontal: 8,
  },
  sendButton: {
    backgroundColor: '#007AFF',
    padding: 10,
    borderRadius: 20,
  },
  sendText: { color: 'white', fontSize: 16 },
  micButton: { fontSize: 24 },
  message: {
    flexDirection: 'row',
    padding: 10,
    marginVertical: 2,
    marginHorizontal: 10,
    borderRadius: 10,
    alignItems: 'flex-start',
  },
  userMsg: { alignSelf: 'flex-end', backgroundColor: '#daf1da' },
  botMsg: { alignSelf: 'flex-start', backgroundColor: '#fff' },
  avatarContainer: { marginRight: 8 },
  userAvatar: {
    width: 30, height: 30, borderRadius: 15,
    backgroundColor: '#007AFF', justifyContent: 'center', alignItems: 'center',
  },
  avatarText: { color: '#fff', fontWeight: 'bold' },
  botAvatar: { width: 30, height: 30, borderRadius: 15 },
  markdown: { flex: 1, lineHeight: 20 },
});
