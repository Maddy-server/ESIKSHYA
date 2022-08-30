-- -- name: CreateScorePoint :exec 
-- INSERT INTO score(
--    user_id,own_points,op_id,op_points,created_at,subject
-- ) VALUES (
--     ?,?,?,?,?,?
-- );

-- -- name: ScoreDetailsList :many
-- SELECT * FROM score WHERE user_id=? ORDER BY created_at DESC;

-- -- name: ScoreDetailsListByCountry :many
-- SELECT DISTINCT score.user_id,children_detail.full_name,children_detail.country,
-- children_detail.grade,  MAX(score.own_points) FROM score  LEFT JOIN  children_detail
-- ON children_detail.children_id = score.user_id WHERE children_detail.country=? GROUP BY score.user_id ORDER BY 5 DESC;

-- -- name: ScoreDetailsListByState :many
-- SELECT DISTINCT score.user_id,children_detail.full_name,children_detail.state,
-- children_detail.grade, MAX(score.own_points) FROM score  LEFT JOIN  children_detail
-- ON children_detail.children_id = score.user_id WHERE children_detail.state=? GROUP BY score.user_id ORDER BY 5 DESC;

-- -- name: ScoreDetailsSats :many
-- SELECT * FROM  score WHERE user_id=? AND created_at >= ?  ORDER BY own_points DESC;