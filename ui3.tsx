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
  SafeAreaView,
  Keyboard,
  TouchableWithoutFeedback
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
  const flatListRef = useRef<FlatList>(null);

  useEffect(() => {
    const loadChat = async () => {
      const history = await AsyncStorage.getItem('chat-history');
      if (history) setMessages(JSON.parse(history));
    };
    loadChat();
  }, []);

  useEffect(() => {
    AsyncStorage.setItem('chat-history', JSON.stringify(messages));
  }, [messages]);

  useEffect(() => {
    Voice.onSpeechResults = e => {
      setInput(e.value?.[0] || '');
    };
    return () => {
      Voice.destroy().then(Voice.removeAllListeners);
    };
  }, []);

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

    const botId = Date.now().toString() + '-bot';
    setMessages(prev => [...prev, { id: botId, text: '...', sender: 'bot' }]);

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

        setMessages(prev =>
          prev.map(msg => (msg.id === botId ? { ...msg, text: botText } : msg))
        );
      }

      Tts.speak(botText);
    } catch (err) {
      console.error('Streaming error:', err);
    }
  };

  const renderItem = ({ item }: { item: Message }) => {
    const isUser = item.sender === 'user';
    return (
      <View
        style={[
          styles.messageRow,
          isUser ? styles.rightAlign : styles.leftAlign,
        ]}
      >
        {!isUser && (
          <Image
            source={require('./assets/bot-avatar.png')}
            style={styles.avatar}
          />
        )}
        <View style={[styles.bubble, isUser ? styles.userBubble : styles.botBubble]}>
          <Markdown style={styles.markdown}>{item.text}</Markdown>
        </View>
        {isUser && (
          <View style={styles.userAvatar}>
            <Text style={styles.userInitials}>U</Text>
          </View>
        )}
      </View>
    );
  };

  return (
    <SafeAreaView style={styles.safeArea}>
      <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
        <KeyboardAvoidingView
          style={styles.container}
          behavior={Platform.OS === 'ios' ? 'padding' : undefined}
          keyboardVerticalOffset={80}
        >
          <FlatList
            ref={flatListRef}
            data={messages}
            keyExtractor={item => item.id}
            renderItem={renderItem}
            onContentSizeChange={() => flatListRef.current?.scrollToEnd({ animated: true })}
            onLayout={() => flatListRef.current?.scrollToEnd({ animated: true })}
            contentContainerStyle={styles.chat}
          />

          <View style={styles.inputBar}>
            <TouchableOpacity onPress={handleVoice}>
              <Text style={styles.icon}>🎤</Text>
            </TouchableOpacity>
            <TextInput
              value={input}
              onChangeText={setInput}
              placeholder="Type your message..."
              style={styles.input}
              multiline
            />
            <TouchableOpacity onPress={sendMessage} style={styles.sendButton}>
              <Text style={styles.sendText}>Send</Text>
            </TouchableOpacity>
          </View>
        </KeyboardAvoidingView>
      </TouchableWithoutFeedback>
    </SafeAreaView>
  );
};

export default App;

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#f9f9f9' },
  container: { flex: 1 },
  chat: { padding: 10 },
  messageRow: {
    flexDirection: 'row',
    marginBottom: 10,
    alignItems: 'flex-end',
  },
  leftAlign: { justifyContent: 'flex-start' },
  rightAlign: { justifyContent: 'flex-end' },
  bubble: {
    maxWidth: '75%',
    padding: 10,
    borderRadius: 20,
  },
  botBubble: {
    backgroundColor: '#e5e5ea',
    marginLeft: 8,
    borderTopLeftRadius: 0,
  },
  userBubble: {
    backgroundColor: '#007AFF',
    marginRight: 8,
    borderTopRightRadius: 0,
  },
  markdown: {
    body: {
      color: '#000',
      fontSize: 16,
    },
    strong: { fontWeight: 'bold' },
    paragraph: { marginTop: 0, marginBottom: 0 },
  },
  avatar: {
    width: 32,
    height: 32,
    borderRadius: 16,
  },
  userAvatar: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: '#007AFF',
    justifyContent: 'center',
    alignItems: 'center',
  },
  userInitials: {
    color: '#fff',
    fontWeight: 'bold',
  },
  inputBar: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 8,
    borderTopWidth: 1,
    borderTopColor: '#ddd',
    backgroundColor: '#fff',
  },
  input: {
    flex: 1,
    minHeight: 40,
    maxHeight: 100,
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 20,
    backgroundColor: '#f0f0f0',
    marginHorizontal: 8,
    fontSize: 16,
  },
  sendButton: {
    backgroundColor: '#007AFF',
    borderRadius: 20,
    paddingHorizontal: 16,
    paddingVertical: 8,
  },
  sendText: {
    color: '#fff',
    fontSize: 16,
  },
  icon: {
    fontSize: 24,
  },
});
