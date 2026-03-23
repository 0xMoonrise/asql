from railroad import Diagram, Sequence, Choice, Optional, OneOrMore, ZeroOrMore, NonTerminal, Terminal

DDL_RULES = {
    "DDL_EXPR": (
        "CREATE TABLE statement",
        Diagram(
            Sequence(
                Terminal("CREATE TABLE"),
                NonTerminal("IDENTIFIER"),
                Terminal("("),
                NonTerminal("TABLE_BODY"),
                Terminal(")")
            )
        )
    ),
    "TABLE_BODY": (
        "Columns and constraints",
        Diagram(
            Sequence(
                Choice(0,
                    NonTerminal("COLUMN_DEF"),
                    NonTerminal("CONSTRAINT_DEF")
                ),
                ZeroOrMore(
                    Sequence(
                        Terminal(","),
                        Choice(0,
                            NonTerminal("COLUMN_DEF"),
                            NonTerminal("CONSTRAINT_DEF")
                        )
                    )
                )
            )
        )
    ),
    "COLUMN_DEF": (
        "Column definition",
        Diagram(
            Sequence(
                NonTerminal("IDENTIFIER"),
                NonTerminal("DATA_TYPE"),
                Optional(NonTerminal("NULLABILITY"))
            )
        )
    ),
    "DATA_TYPE": (
        "Supported data types",
        Diagram(
            Choice(0,
                NonTerminal("NUMERIC_TYPE"),
                NonTerminal("CHAR_TYPE"),
                NonTerminal("DATE"),
            )
        )
    ),
    "NUMERIC_TYPE": (
        "Numeric with precision",
        Diagram(
            Sequence(
                Terminal("NUMERIC"),
                Terminal("("),
                NonTerminal("INTEGER"),
                Optional(Sequence(Terminal(","), NonTerminal("INTEGER"))),
                Terminal(")")
            )
        )
    ),
    "CHAR_TYPE": (
        "Fixed-length string",
        Diagram(
            Sequence(
                Terminal("CHAR"),
                Terminal("("),
                NonTerminal("INTEGER"),
                Terminal(")")
            )
        )
    ),
    "NULLABILITY": (
        "NULL or NOT NULL",
        Diagram(
            Choice(0,
                Sequence(Terminal("NOT"), Terminal("NULL")),
                Terminal("NULL")
            )
        )
    ),
    "CONSTRAINT_DEF": (
        "Named constraint",
        Diagram(
            Sequence(
                Terminal("CONSTRAINT"),
                NonTerminal("IDENTIFIER"),
                NonTerminal("CONSTRAINT_TYPE")
            )
        )
    ),
    "CONSTRAINT_TYPE": (
        "Constraint type",
        Diagram(
            Choice(0,
                Sequence(
                    Terminal("PRIMARY KEY"),
                    Terminal("("), NonTerminal("COL_LIST"), Terminal(")")
                ),
                Sequence(
                    Terminal("FOREIGN KEY"),
                    Terminal("("), NonTerminal("COL_LIST"), Terminal(")"),
                    Terminal("REFERENCES"),
                    NonTerminal("IDENTIFIER"),
                    Terminal("("), NonTerminal("COL_LIST"), Terminal(")")
                ),
                Sequence(
                    Terminal("CHECK"),
                    Terminal("("), NonTerminal("CONDITION"), Terminal(")")
                )
            )
        )
    ),
    "COL_LIST": (
        "One or more columns",
        Diagram(
            OneOrMore(NonTerminal("IDENTIFIER"), Terminal(","))
        )
    ),
}

DML_RULES = {
    "DML_EXPR": (
        "Full SELECT statement",
        Diagram(
            Sequence(
                NonTerminal("SELECT_EXPR"),
                NonTerminal("FROM_EXPR"),
                Optional(NonTerminal("WHERE_CLAUSE"))
            )
        )
    ),
    "SELECT_EXPR": (
        "Column selection",
        Diagram(
            Sequence(
                Terminal("SELECT"),
                Choice(0,
                    Terminal("*"),
                    NonTerminal("COLUMNS_EXPR")
                )
            )
        )
    ),
    "COLUMNS_EXPR": (
        "One or more columns",
        Diagram(
            OneOrMore(NonTerminal("NAME_EXPR"), Terminal(","))
        )
    ),
    "NAME_EXPR": (
        "Simple or qualified name",
        Diagram(
            Choice(0,
                NonTerminal("IDENTIFIER"),
                Sequence(
                    NonTerminal("IDENTIFIER"),
                    Terminal("."),
                    NonTerminal("IDENTIFIER")
                )
            )
        )
    ),
    "FROM_EXPR": (
        "Table sources",
        Diagram(
            Sequence(
                Terminal("FROM"),
                NonTerminal("DATABASES_EXPR")
            )
        )
    ),
    "DATABASES_EXPR": (
        "One or more tables",
        Diagram(
            OneOrMore(NonTerminal("TABLE_EXPR"), Terminal(","))
        )
    ),
    "TABLE_EXPR": (
        "Table or subquery",
        Diagram(
            Choice(0,
                Sequence(
                    NonTerminal("NAME_EXPR"),
                    Optional(NonTerminal("ALIAS"))
                ),
                Sequence(
                    Terminal("("),
                    NonTerminal("DML_EXPR"),
                    Terminal(")"),
                    NonTerminal("ALIAS")
                )
            )
        )
    ),
    "ALIAS": (
        "Table alias",
        Diagram(NonTerminal("IDENTIFIER"))
    ),
    "WHERE_CLAUSE": (
        "Filter condition",
        Diagram(
            Sequence(
                Terminal("WHERE"),
                NonTerminal("CONDITION_EXPR")
            )
        )
    ),
    "CONDITION_EXPR": (
        "Boolean condition",
        Diagram(
            Sequence(
                Optional(Terminal("NOT")),
                NonTerminal("NAME_EXPR"),
                Choice(0,
                    NonTerminal("WHERE_SUBQUERY"),
                    Sequence(
                        NonTerminal("RELATION_EXPR"),
                        Choice(0,
                            NonTerminal("CONSTANT"),
                            NonTerminal("NAME_EXPR")
                        ),
                        Optional(
                            Sequence(
                                Choice(0,
                                    Terminal("AND"),
                                    Terminal("OR")
                                ),
                                Optional(Terminal("NOT")),
                                NonTerminal("CONDITION_EXPR")
                            )
                        )
                    )
                )
            )
        )
    ),
    "CONDITION_EXPR": (
        "Boolean condition",
        Diagram(
            Choice(0,
                Sequence(
                    NonTerminal("NAME_EXPR"),
                    Optional(Terminal("NOT")),
                    NonTerminal("WHERE_SUBQUERY")
                ),
                Sequence(
                    Optional(Terminal("NOT")),
                    NonTerminal("NAME_EXPR"),
                    NonTerminal("RELATION_EXPR"),
                    Choice(0,
                        NonTerminal("CONSTANT"),
                        NonTerminal("NAME_EXPR")
                    ),
                    ZeroOrMore(
                        Sequence(
                            Choice(0,
                                Terminal("AND"),
                                Terminal("OR")
                            ),
                            NonTerminal("CONDITION_EXPR")
                        )
                    )
                )
            )
        )
    ),
    "WHERE_SUBQUERY": (
        "IN subquery with optional continuation",
        Diagram(
            Sequence(
                Terminal("IN"),
                Terminal("("),
                NonTerminal("DML_EXPR"),
                Terminal(")"),
                Optional(
                    Sequence(
                        Choice(0,
                            Terminal("AND"),
                            Terminal("OR")
                        ),
                        Optional(Terminal("NOT")),
                        NonTerminal("CONDITION_EXPR")
                    )
                )
            )
        )
    ),
    "RELATION_EXPR": (
        "Comparison operator",
        Diagram(
            Choice(0,
                Terminal("="),
                Terminal("<"),
                Terminal("<="),
                Terminal(">"),
                Terminal(">="),
                Terminal("<>")
            )
        )
    ),
}
