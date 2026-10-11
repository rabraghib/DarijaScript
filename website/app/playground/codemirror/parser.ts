import { styleTags, tags as t } from '@lezer/highlight';
import { LRParser } from '@lezer/lr';
import {
  foldNodeProp,
  foldInside,
  indentNodeProp,
  delimitedIndent,
} from '@codemirror/language';

const parser = getParser();

export const parserWithMetadata = parser.configure({
  props: [
    styleTags({
      Boolean: t.bool,
      '9ayed': t.definitionKeyword,
      'ilakan sinn ilamakanch ma7dBa9i rjje3 9ayed WAA fonksyon':
        t.controlKeyword,
      'True False s7i7 S7I7 ghalt GHALT': t.bool,
      'CallExpression/FieldExpression/FieldIdentifier': t.function(
        t.propertyName
      ),
      StatementIdentifier: t.labelName,
      Identifier: t.variableName,
      FunctionCall: t.function(t.variableName),
      'FunctionDeclarator/Identifier FunctionDeclarator/DestructorName':
        t.function(t.definition(t.variableName)),
      NamespaceIdentifier: t.namespace,
      OperatorName: t.operator,
      ArithOp: t.arithmeticOperator,
      LogicOp: t.logicOperator,
      BitOp: t.bitwiseOperator,
      CompareOp: t.compareOperator,
      AssignOp: t.definitionOperator,
      UpdateOp: t.updateOperator,
      LineComment: t.lineComment,
      BlockComment: t.blockComment,
      IntegerLiteral: t.number,
      FloatingPointLiteral: t.number,
      StringLiteral: t.string,
      '( )': t.paren,
      '[ ]': t.squareBracket,
      '{ }': t.brace,
      '.': t.derefOperator,
      ', ;': t.separator,
      FunctionDefinition: t.keyword,
      MethodName: t.operator,
    }),
    indentNodeProp.add({
      Application: delimitedIndent({ closing: ')', align: false }),
    }),
    foldNodeProp.add({
      Application: foldInside,
    }),
  ],
});

function getParser() {
  const spec_identifier: Record<string, number> = {
    true: 24,
    false: 24,
    s7i7: 24,
    S7I7: 24,
    ghalt: 24,
    GHALT: 24,
    ilakan: 72,
    sinn: 74,
    ilamakanch: 76,
    ma7dBa9i: 80,
    rjje3: 90,
    '9ayed': 94,
    WAA: 98,
    fonksyon: 106,
  };
  return LRParser.deserialize({
    version: 14,
    states:
      ",lQ]QPOOOOQO'#Do'#DoOOQO'#Cc'#CcOOQO'#Dr'#DrO!jQPO'#ClO$]QQO'#CnO$dQPO'#CmOOQO'#Dq'#DqO#[QQO'#DqO$iQQO'#DmO&UQQO'#CaO!jQPO'#C|OOQO'#Dm'#DmO&]QPO'#DXOOQO'#Dl'#DlOOQO'#Df'#DfQ]QPOOO&dQPO'#DPO&dQPO'#DTO&iQPO'#DYO!bQPO'#D[O&pQPO'#D^O!bQPO'#DbO&{QQO,59WO'SQPO'#CmO!jQPO,58|O'vQPO'#CoOOQO,59X,59XOOQO,5:P,5:PO!jQPO,59POOQO'#Cu'#CuOOQO'#Cw'#CwO!jQPO,59_O!jQPO,59_O!jQPO,59_O!jQPO,59_O!jQPO,59_OOQO,58{,58{OOQO,59h,59hOOQO,59s,59sO(QQPO,59sOOQO-E7d-E7dO]QPO,59kO]QPO,59oO(XQQO,59tOOQO,59t,59tO(`QPO,59vOOQO'#Da'#DaO(eQPO,59xO(jQPO,59zO(oQPO,59|OOQO1G.r1G.rO(tQPO'#CoO({QQO1G.hO)]QQO'#DoO*aQQO'#CcO+eQQO'#DsO+oQPO,59ZO+tQPO'#DwO+|QPO,5:OO,RQQO1G.kOOQO1G.y1G.yO-QQQO1G.yO-yQQO1G.yO.WQQO1G.yO._QQO1G.yOOQO1G/_1G/_O/YQPO1G/VOOQO1G/Z1G/ZOOQO1G/`1G/`OOQO1G/b1G/bOOQO1G/d1G/dOOQO1G/f1G/fO&sQPO'#DdO]QPO1G/hO!jQPO'#DgO0pQPO,5:_OOQO1G.u1G.uO&sQPO'#DhO0xQPO,5:cOOQO1G/j1G/jOOQO7+$V7+$VO]QPO7+$qOOQO7+%S7+%SO1QQQO,5:ROOQO-E7e-E7eOOQO,5:S,5:SOOQO-E7f-E7fOOQO<<H]<<H]",
    stateData:
      "1c~OPOSQOS!_OS~OS^OYROZRO[RO]RO_SOnZOqZOrZOtaOvaOxbOz]O}cO!PdO!ReO!VfO!bQO!dPO~OYROZRO[RO]RO_SOnZOqZOrZO!bQO!dPO~OfiOS!eXe!eXh!eXj!eXl!eXm!eXn!eXo!eX!h!eX!i!eX!j!eX^!eXd!eXW!eX~O_bX~P#[O_jO~OemOS!aXh!aXj!aXl!aXm!aXn!aXo!aX!h!aX!i!aX!j!aX^!aXd!aXW!aX~OhpOjqOlrOmrOnsOotO!hnO!inO!joO~OSuO~P%gOywO~P]O_SO~OS}O~P!jO]!QO!b!PO!d!PO~O^!TO~P%gO_!UO~OYROZRO[RO]RO_SOnZOqZOrZO^!gP~O!b!XO!d!WO~P'XOy!dO~P]OS!gO~P%gOS!hO~OS!iO~OS!jO~O_!kO~O^!gP~P!jOSUi^UidUiWUi~P%gO^!TX^!cX_!cXd!TXd!cXe!cXf!cXh!cXj!cXl!cXm!cXn!cXo!cX!h!cX!i!cX!j!cX~O^VX^!TX_VXdVXd!TXeVXfVXhVXjVXlVXmVXnVXoVX!hVX!iVX!jVX~Od!mO^!gX~P%gO^!oO~Od!pO^!kX~O^!rO~OW!sO~P%gOhpOotO!hnO!inOSgijgilgimgi!jgi^gidgiWgi~OnsO~P,YOhpOotO!hnO!inOSgilgimgi^gidgiWgi~OjqOnsO!joO~P-XOngi~P,YOhpO!hnO!inOSgijgilgimgingiogi!jgi^gidgiWgi~Ou!tOSsiYsiZsi[si]si_sinsiqsirsitsivsixsizsi}si!Psi!Rsi!Vsi!]si!bsi!dsiysi~Od!mO^!ga~Od!pO^!ka~O^!Zad!Za~P%gO!bQZPoY~",
    goto: "*X!lPPPPP!m!w#`P#}PPPPPP$f%T%l&VPPP!wP&ZP&hPPPP!wPP!mPPP!mPPP!m!mP!mP!mP!m&r!mP&|!m'S'^'dPPP'j'|P(}P)h%T*PPPP*Ta^O]`xz{!l!t{[OSZ]`cijmpqrstxz{!U!l!m!tzTOSZ]`cijmpqrstxz{!U!l!m!tQ!OdR!Sf{WOSZ]`cijmpqrstxz{!U!l!m!tzVOSZ]`cijmpqrstxz{!U!l!m!tQzaR{b{VOSZ]`cijmpqrstxz{!U!l!m!t`UO]`xz{!l!tkhSZcijmpqrst!U!mTkUhgpYg|!V!Y!^!`!a!b!c!vaqYg|!V!Y!^!a!vQ!ReS![j!kR!x!pQlUR!l!SQ`OQx]Ty`xQ!n!YR!w!nQ!q![R!y!qW_O]`xQ!ezQ!f{Q!u!lR!z!t`YO]`xz{!l!tQgSQvZQ|cQ!ViS!Yj!UQ!^mQ!_pQ!`qQ!arQ!bsQ!ctR!v!m!PQOSZ]`cdfijmpqrstxz{!U!l!m!t{XOSZ]`cijmpqrstxz{!U!l!m!tT!Zj!UT!]j!k",
    nodeNames:
      '⚠ LineComment BlockComment Program ; ExpressionStatement AssignmentExpression Identifier ] ArrayAccess IntegerLiteral FloatingPointLiteral BooleanLiteral StringLiteral ) ( ParenthesizedExpression MethodInvocation MethodName ArgumentList , [ AssignOp BinaryExpression CompareOp CompareOp LogicOp BitOp BitOp LogicOp ArithOp ArithOp UnaryExpression LogicOp BitOp IfStatement ilakan sinn ilamakanch ForStatement ma7dBa9i } { Block ReturnStatement rjje3 VariableDeclaration 9ayed ThrowStatement WAA Throws Definition FunctionDefinition fonksyon InferredParameters FunctionCall',
    maxTerm: 73,
    nodeProps: [
      ['isolate', -3, 1, 2, 13, ''],
      [
        'group',
        -11,
        4,
        5,
        35,
        39,
        43,
        44,
        46,
        48,
        50,
        52,
        55,
        'Statement',
        -11,
        6,
        7,
        9,
        10,
        11,
        12,
        13,
        16,
        17,
        23,
        32,
        'Expression',
      ],
      ['openedBy', 14, '(', 41, '{'],
      ['closedBy', 15, ')', 42, '}'],
    ],
    skippedNodes: [0, 1, 2],
    repeatNodeCount: 3,
    tokenData:
      "-U~RwX^#lpq#lqr$ars$nst%xtu&avw&uxy'Syz'Xz{'^{|'c|}'h}!O'c!O!P'm!P!Q(X!Q!Z)k!Z![*[!]!^+U!^!_+Z!_!`+c!`!a+k!c!}+s!}#O,X#P#Q,^#Q#R,c#R#S&a#T#o&a#o#p,h#p#q,m#q#r,z#r#s-P#y#z#l$f$g#l#BY#BZ#l$IS$I_#l$I|$JO#l$JT$JU#l$KV$KW#l&FU&FV#l~#qY!_~X^#lpq#l#y#z#l$f$g#l#BY#BZ#l$IS$I_#l$I|$JO#l$JT$JU#l$KV$KW#l&FU&FV#lR$fPqP!_!`$iQ$nOhQ~$qWOY$nZr$nrs%Zs#O$n#O#P%`#P;'S$n;'S;=`%r<%lO$n~%`O]~~%cTOY$nYZ$nZ;'S$n;'S;=`%r<%lO$n~%uP;=`<%l$n~%}SP~OY%xZ;'S%x;'S;=`&Z<%lO%x~&^P;=`<%l%x~&fT!b~tu&a!Q![&a!c!}&a#R#S&a#T#o&a~&zP!j~vw&}~'SOj~~'XO_~~'^O^~~'cOo~~'hOn~~'mOd~~'pP!Q!['s~'xQZ~!Q!['s#R#S(O~(RQ!Q!['s#R#S(O~(^Qo~z{(d!P!Q%x~(gTOz(dz{(v{;'S(d;'S;=`)e<%lO(d~(yVOz(dz{(v{!P(d!P!Q)`!Q;'S(d;'S;=`)e<%lO(d~)eOQ~~)hP;=`<%l(d~)pRY~!O!P)y!Q![)k#R#S*R~*OPZ~!Q!['s~*UQ!Q![)k#R#S*R~*aSY~!O!P)y!Q![)k#R#S*R#T#U*m~*pP#m#n*s~*vP#X#Y*y~*|P#W#X+P~+UO!b~~+ZOS~~+`P!h~!_!`$i~+hPf~!_!`$i~+pP!i~!_!`$i~+xT!d~tu+s!Q![+s!c!}+s#R#S+s#T#o+s~,^Oe~~,cOW~~,hOl~~,mOz~~,rPl~#p#q,u~,zOm~~-POy~~-UOr~",
    tokenizers: [0, 1],
    topRules: { Program: [0, 3] },
    dynamicPrecedences: { '65': -1 },
    specialized: [
      { term: 64, get: (value: string) => Object.hasOwn(spec_identifier, value) ? spec_identifier[value] : -1 },
    ],
    tokenPrec: 748,
  });
}
