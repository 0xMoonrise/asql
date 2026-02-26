"""Frontend Asql."""
import streamlit as st
from code_editor import code_editor
import requests
import pandas as pd
import os
import utils

app_url = os.getenv("APP_URL", "http://localhost:8080")

custom_buttons = [
    {
        "name": "Run",
        "feather": "Play",
        "primary": True,
        "hasText": True,
        "showWithIcon": True,
        "commands": ["submit"],
        "style": {"bottom": "0.44rem", "right": "0.4rem"},
    },
    {
        "name": "Copy",
        "feather": "Copy",
        "hasText": True,
        "commands": ["copyAll"],
        "style": {"bottom": "3.2rem", "right": "0.4rem"},
    }
]

editor_options = {
    "showLineNumbers": True,
    "showGutter": True,
    "wrap": True,
    "fontSize": 14,
}

st.title("Asql")

st.markdown("""
    <style>
        .stTabs [data-baseweb="tab-list"] button p {
            font-size: 1.2rem;
        }
    </style>
""", unsafe_allow_html=True)

response = code_editor(
    "",
    lang="sql",
    height=[10, 30],
    buttons=custom_buttons,
    options=editor_options,
    key="editor"
)

if response["type"] == "submit" and response["text"]:
    query = response["text"]

    st.subheader("Status")

    try:
        res = requests.post(
            f"{app_url}/run",
            data={"text": query}
        )
        res.raise_for_status()
        data = res.json()

        st.session_state["lexer_data"] = data

        if "Errors" in data and data["Errors"]:
            for err in data["Errors"]:
                st.error(f"{err}")
        else:
            st.success("Query executed successfully.")

    except requests.exceptions.RequestException as e:
        st.error(f"Connection error: {e}")
    except Exception as e:
        st.error(f"Error: {e}")


tab_lexer, tab_parser, tab_semantic = st.tabs(["Lexer", "Parser", "Semantic"])
with tab_lexer:
    if "lexer_data" not in st.session_state \
       or not st.session_state["lexer_data"]:
        st.info("Run a query in the Editor tab to see the Lexer results here.")
    else:
        data = st.session_state["lexer_data"]

        st.subheader("Global Table")
        if "GlobalHeaders" in data and "GlobalTable" in data:
            df_global = pd.DataFrame(data["GlobalTable"],
                                     columns=data["GlobalHeaders"])
            for col in ["No.", "Type", "Code"]:
                if col in df_global.columns:
                    df_global[col] = utils.safe_to_numeric(df_global[col])
            df_global.set_index('No.', inplace=True)
            st.dataframe(df_global, width='stretch')

        st.subheader("Identifiers")
        if "IdentifierHeaders" in data and "Identifiers" in data:
            df_ids = pd.DataFrame(data["Identifiers"],
                                  columns=data["IdentifierHeaders"])
            for col in ["Value"]:
                if col in df_ids.columns:
                    df_ids[col] = utils.safe_to_numeric(df_ids[col])
            df_ids.index += 1
            st.dataframe(df_ids, width='stretch')

        st.subheader("Constants")
        if "ConstantHeaders" in data and "Constants" in data:
            df_const = pd.DataFrame(data["Constants"],
                                    columns=data["ConstantHeaders"])
            for col in ["Value"]:
                if col in df_const.columns:
                    df_const[col] = utils.safe_to_numeric(df_const[col])
            df_const.index += 1
            st.dataframe(df_const, width='stretch')
