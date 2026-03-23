"""Utils module."""
import pandas as pd
import railroad
import io
import streamlit as st
import streamlit.components.v1 as components

railroad.VS = 20
railroad.AR = 15

def safe_to_numeric(series):
    """Cast numeric serie."""
    try:
        return pd.to_numeric(series)
    except (ValueError, TypeError):
        return series

def to_svg(diagram) -> str:
    railroad.DEFAULT_STYLE = ""
    out = io.StringIO()
    diagram.writeSvg(out.write)
    svg = out.getvalue()
    custom_style = """<style>
        svg.railroad-diagram { background-color: transparent; }
        svg.railroad-diagram path { stroke-width: 1.5; stroke: #c0c0c0; fill: none; }
        svg.railroad-diagram text { font: bold 13px monospace; fill: #ffffff; text-anchor: middle; }
        svg.railroad-diagram text.label { text-anchor: start; fill: #aaaaaa; }
        svg.railroad-diagram text.comment { font: italic 11px monospace; fill: #aaaaaa; }
        svg.railroad-diagram rect { stroke-width: 1.5; stroke: #c0c0c0; fill: #2d2d4e; }
        svg.railroad-diagram rect.group-box { stroke: #666; stroke-dasharray: 4 4; fill: none; }
    </style>"""
    svg = svg.replace("<style>", "<!-- replaced -->", 1)
    svg = svg.replace("</style>", "-->", 1)
    svg = svg.replace("<svg ", f"{custom_style}<svg ", 1)
    return svg

def grammar_tab(rules: dict):
    for rule, (desc, diagram) in rules.items():
        with st.expander(f"**{rule}** — {desc}"):
            svg = to_svg(diagram)
            components.html(
                f"""
                <div style="background:#1a1a2e; padding:20px; border-radius:8px; display:flex; justify-content:center; align-items:center; min-height:200px;">
                    {svg}
                </div>
                """,
                height=300,
                scrolling=False
            )
