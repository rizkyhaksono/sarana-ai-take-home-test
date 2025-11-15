/**
 * Crypto utilities for encrypting/decrypting auth data
 */

const ENCRYPTION_KEY = process.env.NEXT_PUBLIC_ENCRYPTION_KEY || 'sarana-notes-secret-key-2024';

/**
 * Encrypt data using AES-like encryption (simple XOR with Base64)
 * For production, use proper encryption like crypto-js or Web Crypto API
 */
export function encrypt(text: string): string {
  if (!text) return '';

  try {
    // Simple XOR encryption with Base64 encoding
    const key = ENCRYPTION_KEY;
    let encrypted = '';

    for (let i = 0; i < text.length; i++) {
      const charCode = text.charCodeAt(i) ^ key.charCodeAt(i % key.length);
      encrypted += String.fromCharCode(charCode);
    }

    // Convert to Base64
    return typeof window !== 'undefined'
      ? btoa(encrypted)
      : Buffer.from(encrypted).toString('base64');
  } catch (error) {
    console.error('Encryption error:', error);
    return text;
  }
}

/**
 * Decrypt data
 */
export function decrypt(encryptedText: string): string {
  if (!encryptedText) return '';

  try {
    // Decode from Base64
    const decoded = typeof window !== 'undefined'
      ? atob(encryptedText)
      : Buffer.from(encryptedText, 'base64').toString();

    // XOR decryption
    const key = ENCRYPTION_KEY;
    let decrypted = '';

    for (let i = 0; i < decoded.length; i++) {
      const charCode = decoded.charCodeAt(i) ^ key.charCodeAt(i % key.length);
      decrypted += String.fromCharCode(charCode);
    }

    return decrypted;
  } catch (error) {
    console.error('Decryption error:', error);
    return encryptedText;
  }
}

/**
 * Encrypt and stringify JSON object
 */
export function encryptJSON(obj: unknown): string {
  try {
    const jsonString = JSON.stringify(obj);
    return encrypt(jsonString);
  } catch (error) {
    console.error('JSON encryption error:', error);
    return '';
  }
}

/**
 * Decrypt and parse JSON object
 */
export function decryptJSON<T = unknown>(encryptedText: string): T | null {
  try {
    const decrypted = decrypt(encryptedText);
    return JSON.parse(decrypted) as T;
  } catch (error) {
    console.error('JSON decryption error:', error);
    return null;
  }
}
