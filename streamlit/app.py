"""Frontend Asql."""
import streamlit as st
from code_editor import code_editor
import requests
import pandas as pd
import os
import utils
import logging
import diagrams

# https://wiki.documentfoundation.org/Documentation/SyntaxDiagrams
# https://www.cs.cornell.edu/courses/cs4120/2026sp/notes/

st.set_page_config(page_title="Asql")

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)

app_url = os.getenv("APP_URL", "http://localhost:8080")
custom_buttons = [
    {
        "name": "Parse",
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
    },
    {
        "name": "Exec",
        "feather": "Play",
        "primary": True,
        "hasText": True,
        "commands": ["submit"],
        "style": {"bottom": "0.44rem", "right": "0.4rem"},
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
    height=[10, 20], # type: ignore[arg-type]
    buttons=custom_buttons,
    options=editor_options,
    key="editor"
)

if response["type"] == "submit" and response["text"]:
    query = response["text"]
    st.subheader("Status")

    try:
        res = requests.post(f"{app_url}/run", data={"text": query})
        res.raise_for_status()
        data = res.json()
        st.session_state["data"] = data

        if "Errors" in data and data["Errors"]:
            for err in data["Errors"]:
                st.error(f"{err}")
        elif "parser" in data:
            st.error("Syntax error detected.")
        else:
            st.success("Query executed successfully.")

    except requests.exceptions.RequestException as e:
        logging.error(e)
        st.error("Connection error")
    except Exception as e:
        logging.error(e)
        st.error("Application error")

tab_lexer, tab_parser, tab_AST, tab_grammar = st.tabs(["Lexer", "Parser", "Grammar", "SQL"])
with tab_lexer:
    if "data" not in st.session_state \
       or not st.session_state["data"]:
        st.info("Run a query in the Editor tab to see the Lexer results here.")
    else:
        data = st.session_state["data"]
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

with tab_parser:
    if "data" not in st.session_state or not st.session_state["data"]:
        st.info("Run a query in the Editor tab to see the Parser results here.")
    else:
        data = st.session_state["data"]

        if "Errors" in data and data["Errors"]:
            st.info("A lexer error was detected, please correct it to proceed to parser analysis.")

        elif "parser" in data:
            st.subheader("Parser Error")
            line_num = data.get("parser_line_number", "?")
            line_src = data.get("parser_line", "")
            st.error(data["parser"])
            st.code(f"Line {line_num}: {line_src}", language="sql")

        else:
            st.success("Query is grammatically correct.")

with tab_grammar:

    st.title("SQL Grammar")

    tab_ddl, tab_dml = st.tabs(["DDL — CREATE TABLE", "DML — SELECT"])

    with tab_ddl:
        st.code("""
DDL_EXPR        := CREATE TABLE IDENTIFIER LPAR TABLE_BODY RPAR
TABLE_BODY      := (COLUMN_DEF | CONSTRAINT_DEF)
                   { , (COLUMN_DEF | CONSTRAINT_DEF) }
COLUMN_DEF      := IDENTIFIER DATA_TYPE [ NULLABILITY ]
DATA_TYPE       := NUMERIC_TYPE | CHAR_TYPE | DATE
NUMERIC_TYPE    := NUMERIC LPAR INTEGER [ , INTEGER ] RPAR
CHAR_TYPE       := CHAR LPAR INTEGER RPAR
NULLABILITY     := NOT NULL | NULL
CONSTRAINT_DEF  := CONSTRAINT IDENTIFIER CONSTRAINT_TYPE
CONSTRAINT_TYPE := PRIMARY KEY LPAR COL_LIST RPAR
                 | CHECK LPAR CONDITION RPAR
                 | FOREIGN KEY IDENTIFIER LPAR COL_LIST RPAR
                   REFERENCES IDENTIFIER LPAR COL_LIST RPAR
COL_LIST        := IDENTIFIER { , IDENTIFIER }
        """, language="bnf")
        utils.grammar_tab(diagrams.DDL_RULES)

    with tab_dml:
        st.code("""
DML_EXPR        := SELECT_EXPR FROM_EXPR [ WHERE_CLAUSE ]
SELECT_EXPR     := SELECT * | SELECT COLUMNS_EXPR
COLUMNS_EXPR    := NAME_EXPR { , NAME_EXPR }
NAME_EXPR       := IDENTIFIER | IDENTIFIER . IDENTIFIER
FROM_EXPR       := FROM DATABASES_EXPR
DATABASES_EXPR  := TABLE_EXPR { , TABLE_EXPR }
TABLE_EXPR      := NAME_EXPR [ ALIAS ] | LPAR DML_EXPR RPAR ALIAS
ALIAS           := IDENTIFIER

WHERE_CLAUSE    := WHERE CONDITION_EXPR
CONDITION_EXPR  := NAME_EXPR [NOT] WHERE_SUBQUERY
                 | [NOT] NAME_EXPR RELATION_EXPR
                   (CONSTANT | NAME_EXPR) { (AND | OR) CONDITION_EXPR }
WHERE_SUBQUERY  := IN LPAR DML_EXPR RPAR { (AND | OR) CONDITION_EXPR }
RELATION_EXPR   := = | < | <= | > | >= | <>
        """, language="bnf")
        utils.grammar_tab(diagrams.DML_RULES)
