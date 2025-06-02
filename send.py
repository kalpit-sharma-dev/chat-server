# send_money_agent.py

import streamlit as st
import re
import requests

# --- Replace these with your actual backend endpoints ---
API_ENDPOINTS = {
    "upi": "http://localhost:8000/api/send/upi",
    "phone": "http://localhost:8000/api/send/phone",
    "account": "http://localhost:8000/api/send/account"
}

# --- Extract payment info from input ---
def parse_input(text):
    result = {}
    text = text.lower()

    # Extract amount
    amount_match = re.search(r"(?:rs\.?|₹)\s?(\d+)", text)
    if amount_match:
        result["amount"] = int(amount_match.group(1))

    # UPI ID detection
    upi_match = re.search(r"\b[\w.-]+@[\w.-]+\b", text)
    if upi_match:
        result["method"] = "upi"
        result["address"] = upi_match.group()
        return result if "amount" in result else None

    # Phone number detection
    phone_match = re.search(r"\b\d{10}\b", text)
    if phone_match:
        result["method"] = "phone"
        result["address"] = phone_match.group()
        return result if "amount" in result else None

    # Account + IFSC detection
    acc_match = re.search(r"\b\d{9,18}\b", text)
    ifsc_match = re.search(r"\b[A-Z]{4}0[A-Z0-9]{6}\b", text)
    if acc_match and ifsc_match:
        result["method"] = "account"
        result["account"] = acc_match.group()
        result["ifsc"] = ifsc_match.group()
        return result if "amount" in result else None

    return None

# --- Call backend API based on method ---
def send_money(parsed):
    method = parsed["method"]
    amount = parsed["amount"]

    try:
        if method == "upi":
            payload = {
                "amount": amount,
                "upi_id": parsed["address"]
            }
        elif method == "phone":
            payload = {
                "amount": amount,
                "phone_number": parsed["address"]
            }
        elif method == "account":
            payload = {
                "amount": amount,
                "account_number": parsed["account"],
                "ifsc": parsed["ifsc"]
            }
        else:
            return {"error": "Unsupported method."}

        response = requests.post(API_ENDPOINTS[method], json=payload)
        return response.json()
    except Exception as e:
        return {"error": str(e)}

# --- Streamlit UI ---
st.set_page_config(page_title="Smart Send Money Agent", page_icon="💸")
st.title("💸 AI-Powered Money Transfer")

user_input = st.text_input("Type your command:", placeholder="e.g. Send ₹500 to ankit@upi or Pay 1000 to 9876543210")

if st.button("Send Money"):
    if not user_input.strip():
        st.warning("Please enter a command.")
    else:
        parsed = parse_input(user_input)
        if not parsed:
            st.error("❌ Could not understand the command. Try using UPI, phone, or account+IFSC.")
        else:
            st.write("🔍 Parsed Data:")
            st.json(parsed)

            st.write("⏳ Sending money...")
            result = send_money(parsed)
            st.success("✅ Response:")
            st.json(result)
