-- Seed de Escenarios para SmartStocks
-- Ejecutar: mysql -u root -p smartstocks < database/seed_scenarios.sql

-- ===========================================
-- ESCENARIOS FÁCILES (easy)
-- ===========================================

INSERT INTO simulator_scenarios (id, difficulty, news_content, chart_data, correct_decision, explanation, created_at, expires_at, is_active) VALUES
(UUID(), 'easy',
'Apple reporta ganancias récord: Las ventas del iPhone superaron todas las expectativas con un aumento del 25% respecto al trimestre anterior. El CEO Tim Cook anuncia nuevos productos innovadores para el próximo año.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[145.50,148.20,152.80,155.40,160.20,168.50],"ticker":"AAPL","asset_name":"Apple Inc."}',
'buy',
'Cuando una empresa reporta ganancias récord y supera expectativas, generalmente es una señal positiva. El aumento del 25% en ventas indica fuerte demanda, lo que típicamente impulsa el precio de la acción al alza.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'easy',
'Crisis en Tesla: Se reportan fallas masivas en los frenos de varios modelos. La empresa enfrenta una investigación federal y posibles retiros del mercado de miles de vehículos.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[280.00,275.50,268.30,255.80,240.20,225.60],"ticker":"TSLA","asset_name":"Tesla Inc."}',
'sell',
'Las fallas de seguridad en vehículos son extremadamente negativas para una automotriz. Las investigaciones federales y retiros masivos implican costos enormes y daño reputacional, lo que típicamente causa caídas significativas.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'easy',
'Microsoft mantiene proyecciones estables: La empresa confirma que sus resultados del trimestre están en línea con lo esperado, sin sorpresas positivas ni negativas.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[310.00,312.50,309.80,311.20,313.40,312.00],"ticker":"MSFT","asset_name":"Microsoft Corp."}',
'hold',
'Cuando una empresa cumple exactamente con las expectativas sin sorpresas, el precio tiende a mantenerse estable. No hay catalizador para comprar ni razón para vender, por lo que mantener es la mejor opción.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'easy',
'Amazon anuncia expansión masiva: El gigante del comercio electrónico invertirá $10 mil millones en nuevos centros de distribución en Latinoamérica, creando 50,000 empleos.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[125.40,128.60,132.80,138.20,145.50,152.30],"ticker":"AMZN","asset_name":"Amazon.com Inc."}',
'buy',
'Una expansión significativa en nuevos mercados indica crecimiento futuro. La inversión de $10 mil millones demuestra confianza de la empresa en su capacidad de expandirse, lo cual es positivo para los accionistas.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'easy',
'Netflix pierde 2 millones de suscriptores: Por segundo trimestre consecutivo, la plataforma de streaming reporta pérdida de usuarios en mercados clave como Estados Unidos y Europa.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[380.00,365.40,342.80,318.60,295.20,270.50],"ticker":"NFLX","asset_name":"Netflix Inc."}',
'sell',
'La pérdida continua de suscriptores es muy negativa para un modelo de negocio basado en suscripciones. Dos trimestres seguidos de pérdidas indican un problema estructural, no temporal.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE);

-- ===========================================
-- ESCENARIOS MEDIOS (medium)
-- ===========================================

INSERT INTO simulator_scenarios (id, difficulty, news_content, chart_data, correct_decision, explanation, created_at, expires_at, is_active) VALUES
(UUID(), 'medium',
'Google enfrenta demanda antimonopolio: El Departamento de Justicia de EE.UU. inicia un juicio histórico contra Alphabet, pero analistas creen que el proceso tardará años en resolverse.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[135.20,138.50,136.80,140.20,137.60,139.80],"ticker":"GOOGL","asset_name":"Alphabet Inc."}',
'hold',
'Aunque las demandas antimonopolio suenan graves, cuando los analistas predicen que tardarán años en resolverse, el impacto inmediato en el precio suele ser limitado. La incertidumbre a largo plazo sugiere mantener y observar.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'medium',
'Nvidia supera expectativas pero advierte sobre el futuro: La empresa de chips reporta ganancias 40% arriba de lo esperado, pero el CEO advierte que la demanda podría moderarse en los próximos trimestres.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[420.00,445.80,478.50,512.30,498.60,485.20],"ticker":"NVDA","asset_name":"NVIDIA Corp."}',
'hold',
'Esta es una situación mixta: resultados excelentes pero perspectivas cautelosas. El mercado ya incorporó las buenas noticias en el precio, y la advertencia del CEO crea incertidumbre. Mantener es prudente hasta tener más claridad.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'medium',
'Banco Central sube tasas de interés: La Reserva Federal anuncia un aumento de 0.25% en las tasas, pero señala que podría ser el último del ciclo de ajuste monetario.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[58.20,56.80,54.40,52.60,55.30,57.80],"ticker":"XLF","asset_name":"Financial Select Sector"}',
'buy',
'Aunque las subidas de tasas suelen ser negativas para las acciones, el anuncio de que podría ser la última del ciclo es positivo. Los bancos se benefician de tasas más altas, y el fin del ciclo de ajuste reduce la incertidumbre.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'medium',
'Pfizer anuncia nuevo medicamento prometedor: La farmacéutica revela resultados positivos en Fase 2 para un tratamiento contra el Alzheimer, pero aún falta la Fase 3 que históricamente tiene 50% de fracaso.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[38.50,39.80,42.30,45.60,44.20,46.80],"ticker":"PFE","asset_name":"Pfizer Inc."}',
'hold',
'Los resultados de Fase 2 son prometedores pero no definitivos. Con 50% de probabilidad de fracaso en Fase 3, hay tanto riesgo como oportunidad. Mantener permite esperar más datos sin perderse una potencial subida.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'medium',
'Meta lanza competidor de Twitter: La nueva app Threads alcanza 100 millones de usuarios en 5 días, pero las métricas de engagement son significativamente menores a las de Instagram.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[285.40,298.60,312.80,328.50,335.20,342.80],"ticker":"META","asset_name":"Meta Platforms Inc."}',
'buy',
'100 millones de usuarios en 5 días es un récord histórico. Aunque el engagement inicial sea bajo, Meta tiene experiencia en monetización y puede mejorar el producto. El potencial de crecimiento justifica la compra.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE);

-- ===========================================
-- ESCENARIOS DIFÍCILES (hard)
-- ===========================================

INSERT INTO simulator_scenarios (id, difficulty, news_content, chart_data, correct_decision, explanation, created_at, expires_at, is_active) VALUES
(UUID(), 'hard',
'Intel reporta pérdidas pero anuncia reestructuración: La empresa perdió $500 millones este trimestre, pero el nuevo CEO presenta un plan de transformación de 3 años que incluye recortes de costos del 20% y nueva tecnología de chips.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[32.50,29.80,27.40,25.60,28.30,31.20],"ticker":"INTC","asset_name":"Intel Corp."}',
'buy',
'Aunque las pérdidas son negativas, el mercado ya las anticipó (precio cayó de $32 a $25). La reestructuración agresiva y el nuevo liderazgo pueden ser catalizadores positivos. El precio actual ya refleja las malas noticias, creando oportunidad de compra.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'hard',
'Disney+ alcanza rentabilidad pero pierde a Bob Iger: El servicio de streaming finalmente genera ganancias, pero el legendario CEO anuncia su retiro definitivo sin sucesor claro.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[95.20,102.80,108.50,115.30,112.40,108.80],"ticker":"DIS","asset_name":"The Walt Disney Co."}',
'hold',
'Esta es una situación compleja: noticias operativas positivas pero incertidumbre de liderazgo. La rentabilidad de streaming es un hito importante, pero la salida de Iger sin sucesor crea riesgo. Mantener es prudente hasta conocer el nuevo CEO.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'hard',
'YPF descubre mega yacimiento de petróleo: La petrolera argentina anuncia uno de los mayores descubrimientos de la década, pero el gobierno considera nuevos impuestos a las exportaciones de energía.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[12.80,14.20,15.60,17.80,16.40,18.50],"ticker":"YPF","asset_name":"YPF S.A."}',
'hold',
'El descubrimiento es muy positivo para el valor a largo plazo, pero la amenaza de nuevos impuestos puede erosionar las ganancias. En mercados emergentes, el riesgo político es significativo. Esperar claridad regulatoria es prudente.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'hard',
'MercadoLibre reporta crecimiento récord en Brasil pero enfrenta competencia feroz: Las ventas crecieron 45%, pero Amazon y Shopee están ganando participación de mercado con precios agresivos.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[1250.00,1320.50,1285.80,1380.40,1345.20,1420.60],"ticker":"MELI","asset_name":"MercadoLibre Inc."}',
'buy',
'El crecimiento del 45% demuestra que MercadoLibre sigue siendo el líder. Aunque la competencia es real, la empresa tiene ventajas en logística y fintech (Mercado Pago) que son difíciles de replicar. El crecimiento supera las preocupaciones competitivas.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE),

(UUID(), 'hard',
'Globant gana contrato millonario con la FIFA pero pierde cliente clave: La empresa de tecnología firmó un acuerdo de $200 millones con la FIFA, pero perdió un contrato de $150 millones con Disney que no renovó.',
'{"labels":["Ene","Feb","Mar","Abr","May","Jun"],"prices":[185.40,192.80,188.50,195.20,190.60,193.40],"ticker":"GLOB","asset_name":"Globant S.A."}',
'hold',
'El balance neto es ligeramente positivo ($50M), pero la pérdida de Disney sugiere posibles problemas de retención de clientes. La volatilidad del precio indica que el mercado está indeciso. Mantener y monitorear la retención de clientes.',
NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), TRUE);

-- Verificar escenarios creados
SELECT difficulty, COUNT(*) as cantidad
FROM simulator_scenarios
WHERE is_active = TRUE
GROUP BY difficulty;
