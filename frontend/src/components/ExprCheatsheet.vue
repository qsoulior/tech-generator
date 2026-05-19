<script setup lang="ts">
import { NDrawer, NDrawerContent, NTabs, NTabPane, NFlex, NText, NCard, NTag, NAlert, NScrollbar } from "naive-ui"

defineProps<{
  show: boolean
}>()

defineEmits<{
  "update:show": [value: boolean]
}>()

interface FunctionDoc {
  signature: string
  description: string
  example?: string
}

interface Section {
  title: string
  items: FunctionDoc[]
}

const exprConstants: Section = {
  title: "Константы",
  items: [
    { signature: "pi", description: "Число π.", example: "2 * pi * radius" },
    { signature: "e", description: "Основание натурального логарифма.", example: "pow(e, x)" },
    { signature: "g", description: "Ускорение свободного падения, м/с².", example: "mass * g" },
  ],
}

const exprMath: Section = {
  title: "Математика",
  items: [
    { signature: "sqrt(x)", description: "Квадратный корень.", example: "sqrt(144) // 12" },
    { signature: "pow(x, y)", description: "Возведение в степень xʸ.", example: "pow(2, 10) // 1024" },
    { signature: "exp(x)", description: "Экспонента eˣ.", example: "exp(1) // ≈ 2.718" },
    { signature: "log(x)", description: "Натуральный логарифм.", example: "log(e) // 1" },
    { signature: "log10(x)", description: "Десятичный логарифм.", example: "log10(1000) // 3" },
    { signature: "log2(x)", description: "Логарифм по основанию 2.", example: "log2(8) // 3" },
    { signature: "hypot(x, y)", description: "√(x² + y²) — гипотенуза.", example: "hypot(3, 4) // 5" },
    {
      signature: "mod(x, y)",
      description: "Остаток от деления для вещественных чисел.",
      example: "mod(10, 3) // 1",
    },
    { signature: "abs(x)", description: "Модуль числа.", example: "abs(-5) // 5" },
    { signature: "sign(x)", description: "Знак числа: -1, 0 или 1.", example: "sign(-3.14) // -1" },
    { signature: "ceil(x)", description: "Округление вверх до целого.", example: "ceil(1.2) // 2" },
    { signature: "floor(x)", description: "Округление вниз до целого.", example: "floor(1.7) // 1" },
    { signature: "min(a, b)", description: "Минимум из двух значений.", example: "min(5, 8) // 5" },
    { signature: "max(a, b)", description: "Максимум из двух значений.", example: "max(5, 8) // 8" },
  ],
}

const exprTrig: Section = {
  title: "Тригонометрия",
  items: [
    { signature: "sin(x)", description: "Синус, аргумент в радианах.", example: "sin(pi / 2) // 1" },
    { signature: "cos(x)", description: "Косинус, аргумент в радианах.", example: "cos(pi) // -1" },
    { signature: "tan(x)", description: "Тангенс, аргумент в радианах." },
    { signature: "asin(x)", description: "Арксинус, результат в радианах." },
    { signature: "acos(x)", description: "Арккосинус, результат в радианах." },
    { signature: "atan(x)", description: "Арктангенс, результат в радианах." },
    {
      signature: "atan2(y, x)",
      description: "Арктангенс с учётом квадранта.",
      example: "atan2(1, 1) // ≈ 0.785 (π/4)",
    },
  ],
}

const exprFormat: Section = {
  title: "Округление и форматирование",
  items: [
    { signature: "round(x)", description: "Округление до целого.", example: "round(3.7) // 4" },
    {
      signature: "round(x, n)",
      description: "Округление до n знаков после запятой.",
      example: 'round(3.14159, 2) // 3.14',
    },
    {
      signature: "roundStep(x, step)",
      description: "Округление кратно шагу: цены до 100 ₽, масса до 0,5 кг и т.п.",
      example: "roundStep(1234, 100) // 1200",
    },
    {
      signature: "clamp(x, low, high)",
      description: "Ограничение значения отрезком [low; high].",
      example: "clamp(150, 0, 100) // 100",
    },
    {
      signature: "interpolate(x, x0, x1, y0, y1)",
      description: "Линейная интерполяция между двумя точками (x0,y0) и (x1,y1).",
      example: "interpolate(2.5, 0, 5, 0, 100) // 50",
    },
    {
      signature: "formatNumber(x, decimals)",
      description: "Число в виде строки с десятичной запятой и неразрывными пробелами разрядов.",
      example: 'formatNumber(1234567.5, 2) // "1 234 567,50"',
    },
    {
      signature: "formatNumber(x, decimals, decSep, thouSep)",
      description: "Тот же формат, но с собственными разделителями.",
      example: 'formatNumber(1234.5, 1, ".", ",") // "1,234.5"',
    },
    {
      signature: "percent(x)",
      description: "Доля → проценты (без знаков после запятой по умолчанию).",
      example: 'percent(0.157) // "16%"',
    },
    {
      signature: "percent(x, decimals)",
      description: "Проценты с заданной точностью.",
      example: 'percent(0.157, 1) // "15,7%"',
    },
    {
      signature: "scientific(x)",
      description: "Экспоненциальная запись с двумя знаками после запятой.",
      example: 'scientific(123000.0) // "1.23e+05"',
    },
    {
      signature: "scientific(x, decimals)",
      description: "Экспоненциальная запись с заданной точностью.",
      example: 'scientific(123000.0, 0) // "1e+05"',
    },
  ],
}

const exprStrings: Section = {
  title: "Строки",
  items: [
    { signature: "len(s)", description: "Длина строки или массива.", example: 'len("abc") // 3' },
    { signature: "upper(s)", description: "Верхний регистр.", example: 'upper("abc") // "ABC"' },
    { signature: "lower(s)", description: "Нижний регистр.", example: 'lower("ABC") // "abc"' },
    {
      signature: "trim(s)",
      description: "Убирает пробельные символы по краям.",
      example: 'trim("  abc  ") // "abc"',
    },
    {
      signature: "replace(s, old, new)",
      description: "Замена подстрок.",
      example: 'replace("foo bar", " ", "_") // "foo_bar"',
    },
    {
      signature: "split(s, sep)",
      description: "Разбивает строку по разделителю.",
      example: 'split("a,b,c", ",") // ["a","b","c"]',
    },
    { signature: "string(x)", description: "Приводит число к строке.", example: 'string(42) // "42"' },
    { signature: "int(x)", description: "Приводит значение к целому числу.", example: 'int("42") // 42' },
    {
      signature: "float(x)",
      description: "Приводит значение к вещественному числу.",
      example: 'float("3.14") // 3.14',
    },
  ],
}

const tplBasic: Section = {
  title: "Подстановка и поток управления",
  items: [
    { signature: "{{ .var }}", description: "Подставляет значение переменной по её идентификатору." },
    {
      signature: "{{ if .x }}…{{ end }}",
      description: "Условие: блок выводится, если значение «истинно».",
    },
    {
      signature: "{{ if .x }}…{{ else }}…{{ end }}",
      description: "Условие с альтернативной веткой.",
    },
    {
      signature: "{{ range .items }}…{{ . }}…{{ end }}",
      description: "Цикл по коллекции, точка внутри — текущий элемент.",
    },
    { signature: "{{/* комментарий */}}", description: "Комментарий, в результат не попадает." },
    {
      signature: "{{ .x | upper }}",
      description: "Конвейер: значение слева передаётся последним аргументом в функцию справа.",
    },
  ],
}

const tplStrings: Section = {
  title: "Форматирование текста",
  items: [
    {
      signature: '{{ printf "%.2f" .price }}',
      description: "Форматирование по printf-шаблону.",
      example: '.price = 1234.567 → "1234.57"',
    },
    {
      signature: '{{ default "—" .x }}',
      description: "Значение по умолчанию, если переменная пустая.",
      example: '.x = "" → "—"',
    },
    { signature: "{{ upper .name }}", description: "В верхний регистр.", example: '"abc" → "ABC"' },
    { signature: "{{ lower .name }}", description: "В нижний регистр.", example: '"ABC" → "abc"' },
    {
      signature: "{{ title .name }}",
      description: "Заглавная первая буква каждого слова.",
      example: '"foo bar" → "Foo Bar"',
    },
    {
      signature: '{{ replace " " "_" .name }}',
      description: "Замена подстрок.",
      example: '"foo bar" → "foo_bar"',
    },
    {
      signature: "{{ trim .name }}",
      description: "Убирает пробелы по краям.",
      example: '"  abc  " → "abc"',
    },
    {
      signature: "{{ quote .name }}",
      description: "Оборачивает значение в двойные кавычки.",
      example: 'foo → "foo"',
    },
    {
      signature: '{{ repeat 3 "—" }}',
      description: "Повторение строки.",
      example: '→ "———"',
    },
    {
      signature: '{{ cat "сумма:" .price }}',
      description: "Склеивает значения через пробел.",
      example: '.price = 100 → "сумма: 100"',
    },
    {
      signature: "{{ abbrev 10 .name }}",
      description: "Обрезает строку до указанной длины и добавляет многоточие.",
      example: '"длинное название" → "длинное..."',
    },
    {
      signature: "{{ kebabcase .name }}",
      description: "Преобразование к kebab-форме (только для латиницы).",
      example: '"MyVariable" → "my-variable"',
    },
    {
      signature: "{{ snakecase .name }}",
      description: "Преобразование к snake-форме (только для латиницы).",
      example: '"MyVariable" → "my_variable"',
    },
  ],
}

const tplDates: Section = {
  title: "Даты",
  items: [
    {
      signature: "{{ now }}",
      description: "Текущая дата и время на момент генерации документа.",
    },
    {
      signature: '{{ date "02.01.2006" now }}',
      description: "Форматирование даты по Go-эталону: 02 — день, 01 — месяц, 2006 — год.",
      example: '→ "19.05.2026"',
    },
    {
      signature: '{{ dateModify "24h" now }}',
      description: "Сдвигает дату на указанный интервал (24h, -1h, 30m).",
      example: '+24h → завтра в это же время',
    },
  ],
}

const tplMath: Section = {
  title: "Числа и логика",
  items: [
    { signature: "{{ add 2 3 }}", description: "Сложение.", example: "→ 5" },
    { signature: "{{ sub 10 3 }}", description: "Вычитание.", example: "→ 7" },
    { signature: "{{ mul 3 4 }}", description: "Умножение.", example: "→ 12" },
    { signature: "{{ div 10 3 }}", description: "Целочисленное деление.", example: "→ 3" },
    { signature: "{{ mod 10 3 }}", description: "Остаток от деления.", example: "→ 1" },
    { signature: "{{ max 5 8 }}", description: "Максимум из значений.", example: "→ 8" },
    { signature: "{{ min 5 8 }}", description: "Минимум из значений.", example: "→ 5" },
    {
      signature: '{{ ternary "да" "нет" .x }}',
      description: "Тернарный оператор: значение «истина», значение «ложь», условие.",
      example: '.x = true → "да"',
    },
    {
      signature: "{{ empty .x }}",
      description: "Истина, если значение пустое (0, \"\", nil, пустой список).",
      example: '.x = "" → true',
    },
    {
      signature: '{{ contains "foo" .name }}',
      description: "Истина, если подстрока встречается в значении.",
      example: '"food" → true',
    },
    {
      signature: '{{ hasPrefix "Mr." .name }}',
      description: "Истина, если значение начинается с указанного префикса.",
      example: '"Mr. Jones" → true',
    },
    {
      signature: "{{ len .items }}",
      description: "Длина строки или коллекции.",
      example: "[1, 2, 3] → 3",
    },
  ],
}

const tplLists: Section = {
  title: "Списки",
  items: [
    {
      signature: "{{ first .items }}",
      description: "Первый элемент списка.",
      example: "[1, 2, 3] → 1",
    },
    {
      signature: "{{ last .items }}",
      description: "Последний элемент списка.",
      example: "[1, 2, 3] → 3",
    },
    {
      signature: '{{ join ", " .items }}',
      description: "Склеивает элементы через разделитель.",
      example: '["a", "b", "c"] → "a, b, c"',
    },
    {
      signature: '{{ splitList "," .csv }}',
      description: "Разбивает строку по разделителю в список.",
      example: '"a,b,c" → ["a", "b", "c"]',
    },
  ],
}

const exprSections: Section[] = [exprConstants, exprMath, exprTrig, exprFormat, exprStrings]
const tplSections: Section[] = [tplBasic, tplStrings, tplDates, tplMath, tplLists]

const interpolationSyntax = "{{ .имя }}"
</script>

<template>
  <n-drawer :show="show" :width="640" placement="right" @update:show="$emit('update:show', $event)">
    <n-drawer-content title="Справка по функциям" :native-scrollbar="false" closable>
      <n-tabs type="line" animated default-value="expr">
        <n-tab-pane name="expr" tab="Выражения переменных">
          <n-alert type="info" :show-icon="false" style="margin-bottom: 1rem">
            Используется в формулах вычисляемых переменных и в их ограничениях. Константы
            <n-text code>pi</n-text>, <n-text code>e</n-text> и <n-text code>g</n-text>
            можно использовать без скобок.
          </n-alert>
          <n-scrollbar style="max-height: calc(100vh - 220px)">
            <n-flex vertical :size="24">
              <div v-for="section in exprSections" :key="section.title">
                <n-flex align="center" :size="8" style="margin-bottom: 0.5rem">
                  <n-text strong>{{ section.title }}</n-text>
                  <n-tag size="small" :bordered="false">{{ section.items.length }}</n-tag>
                </n-flex>
                <n-flex vertical :size="8">
                  <n-card v-for="item in section.items" :key="item.signature" size="small">
                    <n-flex vertical :size="8">
                      <n-text code strong class="cheat-signature">{{ item.signature }}</n-text>
                      <n-text depth="2" class="cheat-description">{{ item.description }}</n-text>
                      <n-text v-if="item.example" depth="3" class="cheat-example">
                        {{ item.example }}
                      </n-text>
                    </n-flex>
                  </n-card>
                </n-flex>
              </div>
            </n-flex>
          </n-scrollbar>
        </n-tab-pane>
        <n-tab-pane name="tpl" tab="Шаблон документа">
          <n-alert type="info" :show-icon="false" style="margin-bottom: 1rem">
            Используется в теле шаблона документа. Переменные доступны как
            <n-text code>{{ interpolationSyntax }}</n-text>.
          </n-alert>
          <n-scrollbar style="max-height: calc(100vh - 220px)">
            <n-flex vertical :size="24">
              <div v-for="section in tplSections" :key="section.title">
                <n-flex align="center" :size="8" style="margin-bottom: 0.5rem">
                  <n-text strong>{{ section.title }}</n-text>
                  <n-tag size="small" :bordered="false">{{ section.items.length }}</n-tag>
                </n-flex>
                <n-flex vertical :size="8">
                  <n-card v-for="item in section.items" :key="item.signature" size="small">
                    <n-flex vertical :size="8">
                      <n-text code strong class="cheat-signature">{{ item.signature }}</n-text>
                      <n-text depth="2" class="cheat-description">{{ item.description }}</n-text>
                      <n-text v-if="item.example" depth="3" class="cheat-example">
                        {{ item.example }}
                      </n-text>
                    </n-flex>
                  </n-card>
                </n-flex>
              </div>
            </n-flex>
          </n-scrollbar>
        </n-tab-pane>
      </n-tabs>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped>
.cheat-signature {
  align-self: flex-start;
}

.cheat-description {
  font-size: 0.875rem;
  line-height: 1.45;
}

.cheat-example {
  font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
  font-size: 0.78rem;
  font-style: italic;
}
</style>
